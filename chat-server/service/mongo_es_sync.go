package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"log"
	"strings"
)

func StartMongoToEsSync(ctx context.Context, mongoDB *mongo.Database, esClient *elasticsearch.Client, collectionName, esIndex string) error {
	//mongoDB := initialize.MongoDB   // 你在初始化时保存的 *mongo.Database
	//esClient := initialize.EsClient // 你在初始化时保存的 *elasticsearch.Client

	collection := mongoDB.Collection(collectionName)
	pipeline := mongo.Pipeline{} // 可以自定义过滤条件
	cs, err := collection.Watch(ctx, pipeline)
	if err != nil {
		log.Fatalf("监听 Change Stream 失败: %v", err)
	}
	defer cs.Close(ctx)

	fmt.Println("开始监听 MongoDB 变更并同步到 ES...")

	for cs.Next(ctx) {
		var event bson.M
		if err := cs.Decode(&event); err != nil {
			log.Println("解码 Change Stream 事件失败:", err)
			continue
		}
		log.Printf("收到 Change Stream 事件: %v\n", event)
		opType, _ := event["operationType"].(string)
		log.Printf("操作类型 (长度: %d): %q\n", len(opType), opType) // %q 会用引号包围字符串，并转义特殊字符
		opType = strings.TrimSpace(opType)
		log.Printf("操作类型 (清理后，长度: %d): %q\n", len(opType), opType) // 再次打印确认

		doc := event["fullDocument"]
		log.Printf("fullDocument 的实际类型是: %T\n", doc) // <-- 添加这行

		switch opType {
		case "insert", "replace", "update":
			fmt.Println("成功进入 insert、replace、update")

			// 1. 将 fullDocument (bson.D) 编码回 BSON 字节
			docBytes, err := bson.Marshal(doc)
			if err != nil {
				log.Printf("Marshal fullDocument 到 BSON 字节失败: %v\n", err)
				continue
			}

			// 2. 将 BSON 字节解码到 bson.M (确保是 map 类型)
			var docMap bson.M
			err = bson.Unmarshal(docBytes, &docMap)
			if err != nil {
				log.Printf("Unmarshal BSON 字节到 bson.M 失败: %v\n", err)
				continue
			}
			// 获取文档ID
			var id string
			// 尝试断言为 bson.ObjectID
			if mongoID, ok := docMap["_id"].(bson.ObjectID); ok {
				id = mongoID.Hex()                         // 使用 Hex() 方法获取十六进制字符串
				log.Printf("成功将 _id 转换为十六进制字符串: %s\n", id) // 添加日志确认
			} else {
				// 如果 _id 不是 bson.ObjectID 类型，则回退到通用字符串格式化
				id = fmt.Sprintf("%v", docMap["_id"])
				log.Printf("警告: 文档 _id 类型不是 bson.ObjectID，而是 %T. 使用通用字符串格式化: %s\n", docMap["_id"], id)
			}
			// 从 docMap 中删除 _id 字段，以免它被包含在 JSON Body 中
			delete(docMap, "_id")
			// 转成 JSON
			jsonBody, err := json.Marshal(docMap)
			if err != nil {
				log.Printf("Marshal docMap 到 JSON 失败: %v\n", err)
				continue
			}
			if jsonBody == nil || len(jsonBody) == 0 {
				log.Println("Marshal JSON 结果为空，跳过写入ES。原始 docMap:", docMap)
				continue
			}
			// 写入 ES
			fmt.Println("开始写入es" + string(jsonBody))
			res, err := esClient.Index(
				esIndex,
				strings.NewReader(string(jsonBody)),
				esClient.Index.WithDocumentID(id),
				esClient.Index.WithContext(ctx),
			)
			if err != nil {
				log.Println("写入 ES 失败:", err)
			} else {
				if res.IsError() { // <<< 检查这个错误
					log.Printf("ES 索引文档失败: %s - 文档ID: %s, 响应: %s\n", res.Status(), id, res.String())
				} else {
					log.Printf("ES 文档同步成功: %s, 文档ID: %s\n", res.Status(), id) // <<< 成功日志
				}
				res.Body.Close()
			}

		case "delete":
			// 删除 ES 文档
			docKey := event["documentKey"].(bson.M)
			id := fmt.Sprintf("%v", docKey["_id"])
			res, err := esClient.Delete(esIndex, id, esClient.Delete.WithContext(ctx))
			if err != nil {
				log.Println("删除 ES 文档失败:", err)
			} else {
				res.Body.Close()
			}
		default: // 添加 default case，捕获所有未匹配的情况
			log.Printf("警告：opType '%s' (%q) 未匹配任何已知操作类型。原始事件: %v\n", opType, opType, event)
		}
	}

	if err := cs.Err(); err != nil {
		log.Fatalf("Change Stream 监听出错: %v", err)
	}

	return nil
}
