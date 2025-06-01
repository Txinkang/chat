package initialize

import (
	"chat-server/global"
	"context"
	"github.com/elastic/go-elasticsearch/v9"
	"log/slog"
	"net/http"
	"time"
)

func InitElasticSearch() error {
	slog.Info("初始化elasticsearch")
	EsConfig := global.CHAT_CONFIG.ElasticSearch

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
		Transport: &http.Transport{
			MaxIdleConns:          100,              // 保持的最大空闲连接数
			IdleConnTimeout:       90 * time.Second, // 空闲连接超时时间
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // 15秒超时
	defer cancel()

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		slog.Error("Elasticsearch 连接失败: ", "err", err)
		return err
	}

	res, err := client.Info(client.Info.WithContext(ctx))
	if err != nil {
		slog.Error("Elasticsearch Ping 失败: ", "err", err)
		closeErr := res.Body.Close()
		if closeErr != nil {
			slog.Error("Elasticsearch Ping 失败后，，响应数据流关闭失败: ", "closeErr", closeErr)
		}
		return err
	}

	if res.IsError() {
		slog.Error("Elasticsearch 返回错误: ", "err", err)
		return err
	}
	global.CHAT_ES = client
	slog.Info("Elasticsearch连接成功")
	return nil
}
