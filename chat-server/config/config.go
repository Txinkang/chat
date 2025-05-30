package config

type AppConfig struct {
	Server        Server        `mapstructure:"server" yaml:"server"`
	Mysql         Mysql         `mapstructure:"mysql" yaml:"mysql"`
	Mongo         Mongo         `mapstructure:"mongo" yaml:"mongo"`
	ElasticSearch ElasticSearch `mapstructure:"elasticsearch" yaml:"elasticsearch"`
	Redis         Redis         `mapstructure:"redis" yaml:"redis"`
	Logger        Logger        `mapstructure:"logger" yaml:"logger"`
}
