package config

type ElasticSearch struct {
	Address            string `mapstructure:"address" yaml:"address"`
	ApiKey             string `mapstructure:"api_key" yaml:"api_key"`
	InsecureSkipVerify bool   `mapstructure:"insecure_skip_verify" yaml:"insecure_skip_verify"`
}
