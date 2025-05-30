package config

type Mongo struct {
	URI    string `mapstructure:"uri" yaml:"uri"`
	DBName string `mapstructure:"dbname" yaml:"dbname"`
}
