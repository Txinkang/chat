package service

import (
	"chat-server/global"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"log/slog"

	"strings"
)

func StartMongoToEsSync(ctx context.Context, collectionName string, esIndex string) error {
	slog.Info("初始化 Mongo-Es 数据同步流程")
	collection := global.CHAT_MONGODB.Collection(collectionName)
	pipeline := mongo.Pipeline{}

	//启动监听
	watch, err := collection.Watch(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("监听 Change Stream 失败: %v", err)
	}
	defer watch.Close(ctx)
	slog.Info("开始监听 MongoDB 变更并同步到 ES...")

	//开始监听
	for watch.Next(ctx) {
		//解码文档 bson.M
		var event bson.M
		if err := watch.Decode(&event); err != nil {
			slog.Error("解码 Change Stream 事件失败:", "err", err)
			continue
		}
		operationType, ok := event["operationType"].(string)
		if !ok {
			slog.Error("无法获取 operationType 或其类型不是 string", "event", event)
			continue
		}
		operationType = strings.TrimSpace(operationType)
		document, ok := event["fullDocument"].(bson.M)
		if !ok {
			// 1. 将 fullDocument (bson.D) 编码回 BSON 字节
			docBytes, err := bson.Marshal(event["fullDocument"])
			if err != nil {
				slog.Error("document 到 BSON 字节失败:", "err", err)
				continue
			}

			// 2. 将 BSON 字节解码到 bson.M (确保是 map 类型)
			var docMap bson.M
			err = bson.Unmarshal(docBytes, &docMap)
			if err != nil {
				slog.Error("BSON 字节到 bson.M 失败: ", "err", err)
				continue
			}
			document = docMap
		}

		switch operationType {
		case "insert", "replace", "update":
			//先获取id，再把文档里的_id删掉，不然写入es的时候就会多一个_id
			var id string
			if mongoID, ok := document["_id"].(bson.ObjectID); ok {
				id = mongoID.Hex() // 使用 Hex() 方法获取十六进制字符串
			} else {
				id = fmt.Sprintf("%v", document["_id"])
				slog.Warn("警告: 文档 _id 类型不是 bson.ObjectID,使用通用字符串格式化: ", "id", id)
			}
			delete(document, "_id")

			//把文档转换为json，为了写入es
			jsonDocument, err := json.Marshal(document)
			if err != nil {
				slog.Error("document 到 JSON 失败", "err", err)
				continue
			}
			if jsonDocument == nil || len(jsonDocument) == 0 {
				slog.Error("JSON 结果为空，跳过写入ES。原始 document:", document)
				continue
			}

			//写入es
			res, err := global.CHAT_ES.Index(
				esIndex,
				strings.NewReader(string(jsonDocument)),
				global.CHAT_ES.Index.WithDocumentID(id),
				global.CHAT_ES.Index.WithContext(ctx),
			)
			if err != nil {
				slog.Error("写入 ES 失败:", "err", err)
			} else {
				if res.IsError() { // <<< 检查这个错误
					slog.Error("ES 索引文档失败:", "文档", res.Status(), "文档ID:", id, "响应:", res.String())
				} else {
					slog.Info("ES 文档同步成功: ", "文档", res.Status(), "文档ID:", id) // <<< 成功日志
				}
				res.Body.Close()
			}

		case "delete":
			documentKey := event["documentKey"].(bson.M)
			id := fmt.Sprintf("%v", documentKey["_id"])
			res, err := global.CHAT_ES.Delete(esIndex, id, global.CHAT_ES.Delete.WithContext(ctx))
			if err != nil {
				slog.Error("删除 ES 文档失败:", "err", err)
			} else {
				res.Body.Close()
			}

		default:
			slog.Error("未匹配到任何操作类型:", "operationType", operationType)
		}

		if err := watch.Err(); err != nil {
			slog.Error("Change Stream 监听出错：", "err", err)
		}
	}
	return nil
}
