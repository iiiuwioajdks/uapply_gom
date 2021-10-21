package config

type ServerConf struct {
	Port       int        `json:"port"`
	MysqlInfo  MysqlConf  `json:"mysql"`
	RedisInfo  RedisConf  `json:"redis"`
	LoggerInfo LoggerConf `json:"log"`
	JwtInfo    JwtConf    `json:"jwt"`
}

type JwtConf struct {
	SigningKey string `json:"key"`
}

type MysqlConf struct {
	Host     string `json:"host"`
	UserName string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	Port     int    `json:"port"`
}

type RedisConf struct {
	Host string `json:"host"`
	Auth string `json:"auth"`
	Port int    `json:"port"`
	DB   int    `json:"db"`
}

type LoggerConf struct {
	Level    string `json:"level"`
	Filename string `json:"filename"`
	Mode     string `json:"mode"`
	MaxSize  int    `json:"max_size"`
	MaxAge   int    `json:"max_age"`
	BackUp   int    `json:"max_backups"`
}

type NacosConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Pwd         string `mapstructure:"pwd"`
	Group       string `mapstructure:"group"`
	DataId      string `mapstructure:"dataid"`
	NamespaceId string `mapstructure:"namespaceid"`
}
