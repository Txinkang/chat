package initialize

import (
	"context"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"time"
)

var EsClient *elasticsearch.Client

func InitElasticSearch() error {
	EsConfig := AppConfig.ElasticSearch

	cfg := elasticsearch.Config{
		Addresses: []string{
			EsConfig.Address,
		},
		APIKey: EsConfig.ApiKey,
		// 由于是自签名证书或云服务证书，通常需要跳过TLS验证，但在生产环境中不推荐
		// 或者你可以正确配置CA证书
		//Transport: &http.Transport{
		//	TLSClientConfig: &tls.Config{InsecureSkipVerify: EsConfig.InsecureSkipVerify},
		//},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15秒超时
	defer cancel()

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return fmt.Errorf("Elasticsearch 连接失败: %w", err)
	}

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("Elasticsearch Ping 失败: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("Elasticsearch 返回错误: %s", res.String())
	}
	EsClient = client
	fmt.Println("✅ Elasticsearch 连接成功")
	return nil
}
