package conf

import (
	"github.com/spf13/viper"
	"log"
)

var (
	ConfigYamlInstance ConfigYaml
)

func init() {
	configYamlViper := viper.New()
	configYamlViper.SetConfigFile("./config/config.yaml")
	if err := configYamlViper.ReadInConfig(); err != nil {
		log.Fatalln("读配置文件失败", err)
	}
	if err := configYamlViper.Unmarshal(&ConfigYamlInstance); err != nil {
		log.Fatalln("解析配置文件失败", err)
	}
}

type ConfigYaml struct {
	AppConfig      AppConfig      `mapstructure:"app"`
	LogConfig      LogConfig      `mapstructure:"log"`
	HTTPConfig     HttpConfig     `mapstructure:"http"`
	DingtalkConfig DingtalkConfig `mapstructure:"dingtalk"`
}

type AppConfig struct {
	AppName string `json:"app_name" mapstructure:"app_name"`
}
type LogConfig struct {
	LogsDir           string `json:"logs_dir" mapstructure:"logs_dir"`
	LogsRotationTime  int    `json:"logs_rotation_time" mapstructure:"logs_rotation_time"`
	LogsRotationCount int    `json:"logs_rotation_count" mapstructure:"logs_rotation_count"`
}
type HttpConfig struct {
	ServerListenPort int            `json:"server_listen_port" mapstructure:"server_listen_port"`
	Cors             HttpCORSConfig `json:"cors" mapstructure:"cors"`
}
type HttpCORSConfig struct {
	Enable                    bool   `json:"enable" mapstructure:"enable"`
	AccessControlAllowOrigin  string `json:"access_control_allow_origin" mapstructure:"access_control_allow_origin"`
	AccessControlAllowMethods string `json:"access_control_allow_methods" mapstructure:"access_control_allow_methods"`
	AccessControlAllowHeaders string `json:"access_control_allow_headers" mapstructure:"access_control_allow_headers"`
}
type DingtalkConfig struct {
	TimeOut         int    `json:"time_out" mapstructure:"time_out"`
	MessageDuration int    `json:"message_duration" mapstructure:"message_duration"`
	Boundary        string `json:"boundary" mapstructure:"boundary"`
}
