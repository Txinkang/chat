package config

type Redis struct {
	Address  string `mapstructure:"address" yaml:"address"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}
