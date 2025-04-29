package configs

import "github.com/spf13/viper"

type conf struct {
	RedisHost              string `mapstructure:"REDIS_HOST"`
	RedisPort              string `mapstructure:"REDIS_PORT"`
	RateLimitIp            string `mapstructure:"RATE_LIMIT_IP"`
	RateLimitToken         string `mapstructure:"RATE_LIMIT_TOKEN"`
	RateLimitWindowIP      string `mapstructure:"RATE_LIMIT_WINDOW_IP"`
	RateLimitWindowToken   string `mapstructure:"RATE_LIMIT_WINDOW_TOKEN"`
	RateLimitBlockWindowIP string `mapstructure:"RATE_LIMIT_BLOCK_WINDOW_IP"`
	TokensConfigLimit      string `mapstructure:"TOKENS_CONFIG_LIMIT"`
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	GlobalRateLimit        string `mapstructure:"GLOBAL_RATE_LIMIT"`
	AllowIPLimit           string `mapstructure:"ALLOW_IP_LIMIT"`
	AllowTokenLimit        string `mapstructure:"ALLOW_TOKEN_LIMIT"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
