package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	MongoURI   string `mapstructure:"MONGO_URI"`
	JWTSecret  string `mapstructure:"JWT_SECRET"`
	JWTExpires int    `mapstructure:"JWT_EXPIRES"` // 시간 단위: 시간
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// 기본값 설정
	viper.SetDefault("SERVER_PORT", "8001")
	viper.SetDefault("JWT_EXPIRES", 24) // 24시간

	if err := viper.ReadInConfig(); err != nil {
		// .env 파일이 없어도 환경변수로 실행 가능하게
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
