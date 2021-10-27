package config

type ServerConf struct {
	Port       int        `json:"port" mapstructure:"port"`
	MysqlInfo  MysqlConf  `json:"mysql" mapstructure:"mysql"`
	RedisInfo  RedisConf  `json:"redis" mapstructure:"redis"`
	LoggerInfo LoggerConf `json:"log" mapstructure:"log"`
	JwtInfo    JwtConf    `json:"jwt" mapstructure:"jwt"`
}

type JwtConf struct {
	SigningKey string `json:"key" mapstructure:"key"`
}

type MysqlConf struct {
	Host     string `json:"host" mapstructure:"host"`
	UserName string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	DBName   string `json:"dbname" mapstructure:"dbname"`
	Port     int    `json:"port" mapstructure:"port"`
}

type RedisConf struct {
	Host string `json:"host" mapstructure:"host"`
	Auth string `json:"auth" mapstructure:"auth"`
	Port int    `json:"port" mapstructure:"port"`
	DB   int    `json:"db" mapstructure:"db"`
}

type LoggerConf struct {
	Level    string `json:"level" mapstructure:"level"`
	Filename string `json:"filename" mapstructure:"filename"`
	Mode     string `json:"mode" mapstructure:"mode"`
	MaxSize  int    `json:"max_size" mapstructure:"max_size"`
	MaxAge   int    `json:"max_age" mapstructure:"max_age"`
	BackUp   int    `json:"max_backups" mapstructure:"max_backups"`
}
