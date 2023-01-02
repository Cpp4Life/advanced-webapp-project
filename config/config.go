package config

import "github.com/spf13/viper"

type Config struct {
	Port                    string `mapstructure:"PORT"`
	DBHost                  string `mapstructure:"DB_HOST"`
	DBPort                  string `mapstructure:"DB_PORT"`
	DBUser                  string `mapstructure:"DB_USER"`
	DBPass                  string `mapstructure:"DB_PASS"`
	DBName                  string `mapstructure:"DB_NAME"`
	SecretKey               string `mapstructure:"SECRET_KEY"`
	SMTPHost                string `mapstructure:"SMTP_HOST"`
	SMTPUser                string `mapstructure:"SMTP_USER"`
	SMTPPass                string `mapstructure:"SMTP_PASS"`
	SMTPPort                int    `mapstructure:"SMTP_PORT"`
	GoogleOauthClientId     string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleOauthClientSecret string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOauthRedirectUrl  string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
	GoogleOauthAPI          string `mapstructure:"GOOGLE_OAUTH_API"`
	FrontendLoginUrl        string `mapstructure:"FRONTEND_LOGIN_REDIRECT_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
