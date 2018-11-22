package config

import "github.com/spf13/viper"

// RedisConf - Redis configuration
type RedisConf struct {
	Address  string `mapstructure:"address"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"database"`
}

// ServerConf - server configuration
type ServerConf struct {
	Port int `mapstructure:"port"`
}

// BotConf - bot configuration
type BotConf struct {
	Token string `mapstructure:"token"`
	Debug bool   `mapstructure:"debug"`
}

// OauthConf - Google Oauth2 configuration
type OauthConf struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

var (
	// RedisConfig - Redis configuration
	RedisConfig RedisConf

	// ServerConfig - server configuration
	ServerConfig ServerConf

	// BotConfig - Telegram bot configuration
	BotConfig BotConf

	// OauthConfig - Google Oauth2 configuration
	OauthConfig OauthConf
)

// Init - initialize config
func Init() error {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	redisConfig := viper.Sub("redis")
	serverConfig := viper.Sub("server")
	botConfig := viper.Sub("bot")
	oauthConfig := viper.Sub("oauth")

	if err := redisConfig.Unmarshal(&RedisConfig); err != nil {
		return err
	}

	if err := serverConfig.Unmarshal(&ServerConfig); err != nil {
		return err
	}

	if err := botConfig.Unmarshal(&BotConfig); err != nil {
		return err
	}

	if err := oauthConfig.Unmarshal(&OauthConfig); err != nil {
		return err
	}

	return nil
}
