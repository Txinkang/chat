package main

import (
	"bytes" // 新增导入，用于构建 JSON 请求体
	"context"
	"crypto/tls"
	"encoding/json" // 新增导入，使用标准库的 json 包
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
	// "github.com/elastic/go-elasticsearch/v9/esutil" // 移除这行，因为不再使用 esutil 包的那些方法
)

func main() {
	// 1. Elasticsearch 连接配置
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://4ff4cca5c85e41169f6e3fc0ce869900.ap-east-1.aws.elastic-cloud.com:443", // 你的 Elastic Cloud URL
		},
		APIKey: "SndVaUY1Y0JNMVVNVmZZQ1JFd0s6WlYtZkNfS3dMc2hFbVBMZUJmUnhMdw==", // 你的 API Key

		// 由于是自签名证书或云服务证书，通常需要跳过TLS验证，但在生产环境中不推荐
		// 或者你可以正确配置CA证书
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// 2. 创建 Elasticsearch 客户端
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 3. 检查连接状态
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	} else {
		var r map[string]interface{}
		// 替换 esutil.DecodeJSON
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		fmt.Printf("Client: %s\n", elasticsearch.Version)
		fmt.Printf("Server: %s (%s)\n", r["version"].(map[string]interface{})["number"], r["tagline"])
	}

	// 4. 索引文档 (POST /messages/_doc)
	// 示例文档
	doc := map[string]interface{}{
		"room_id":   "room_001",
		"sender_id": "user_001",
		"type":      "text",
		"content": map[string]interface{}{
			"text": "Hello Elasticsearch from Go!",
		},
		"created_at": 1678886400000, // UTC 毫秒时间戳
	}

	// 替换 esutil.NewJSONReader
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Fatalf("Error encoding document: %s", err)
	}

	res, err = es.Index(
		"messages",                   // 索引名称
		&buf,                         // 传入 bytes.Buffer 的指针作为请求体
		es.Index.WithDocumentID("1"), // 可选：指定文档ID
		es.Index.WithRefresh("true"), // 刷新索引，以便立即搜索到
	)
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error indexing document: %s", res.String())
	} else {
		fmt.Println("\n--- Document indexed successfully ---")
		fmt.Println(res.String())
	}

	// 索引一个包含 image 字段的文档
	docImage := map[string]interface{}{
		"room_id":   "room_002",
		"sender_id": "user_002",
		"type":      "image",
		"content": map[string]interface{}{
			"image": map[string]interface{}{
				"url":    "https://example.com/image_001.jpg",
				"name":   "my_beautiful_image.jpg",
				"size":   102400, // bytes
				"format": "jpeg",
			},
		},
		"created_at": 1678886500000,
	}

	// 替换 esutil.NewJSONReader
	var bufImage bytes.Buffer
	if err := json.NewEncoder(&bufImage).Encode(docImage); err != nil {
		log.Fatalf("Error encoding image document: %s", err)
	}

	res, err = es.Index(
		"messages",
		&bufImage, // 传入 bytes.Buffer 的指针作为请求体
		es.Index.WithDocumentID("2"),
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Fatalf("Error indexing image document: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error indexing image document: %s", res.String())
	} else {
		fmt.Println("\n--- Image document indexed successfully ---")
		fmt.Println(res.String())
	}

	// 5. 搜索文档 (GET /messages/_search)
	fmt.Println("\n--- Searching for documents ---")
	searchQuery := `{"query":{"match_all":{}}}` // 匹配所有文档的简单查询

	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("messages"),
		es.Search.WithBody(strings.NewReader(searchQuery)), // Search.WithBody 仍然接受 io.Reader
		es.Search.WithTrackTotalHits(true),                 // 跟踪总命中数
	)
	if err != nil {
		log.Fatalf("Error searching documents: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error searching documents: %s", res.String())
	} else {
		var r map[string]interface{}
		// 替换 esutil.DecodeJSON
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the search response body: %s", err)
		}
		fmt.Printf("Total hits: %v\n", r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"])
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			fmt.Printf("  ID: %s, Source: %s\n", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}
	}

	// 6. 示例：使用 ngram 子字段搜索
	fmt.Println("\n--- Searching using ngram sub-field ---")
	ngramSearchQuery := `{"query":{"match":{"content.text.ngram":"hello"}}}` // 搜索 "hello" 字符串
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("messages"),
		es.Search.WithBody(strings.NewReader(ngramSearchQuery)), // Search.WithBody 仍然接受 io.Reader
		es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error searching with ngram: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error searching with ngram: %s", res.String())
	} else {
		var r map[string]interface{}
		// 替换 esutil.DecodeJSON
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing ngram search response: %s", err)
		}
		fmt.Printf("Total hits (ngram): %v\n", r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"])
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			fmt.Printf("  ID: %s, Source: %s\n", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		}
	}
}
