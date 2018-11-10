package config

import "github.com/spf13/viper"

// DatabaseConf - database configuration
type DatabaseConf struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
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
	// DatabaseConfig - database configuration
	DatabaseConfig DatabaseConf

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

	databaseConfig := viper.Sub("database")
	serverConfig := viper.Sub("server")
	botConfig := viper.Sub("bot")
	oauthConfig := viper.Sub("oauth")

	if err := databaseConfig.Unmarshal(&DatabaseConfig); err != nil {
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
