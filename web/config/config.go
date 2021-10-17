package config

type ServerConf struct {
	Port       int        `mapstructure:"port"`
	MysqlInfo  MysqlConf  `mapstructure:"mysql"`
	RedisInfo  RedisConf  `mapstructure:"redis"`
	LoggerInfo LoggerConf `mapstructure:"log"`
}

type MysqlConf struct {
	Host     string `mapstructure:"host"`
	UserName string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
}

type RedisConf struct {
	Host string `mapstructure:"host"`
	Auth string `mapstructure:"auth"`
	Port int    `mapstructure:"port"`
	DB   int    `mapstructure:"db"`
}

type LoggerConf struct {
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
	Mode     string `mapstructure:"mode"`
	MaxSize  int    `mapstructure:"max_size"`
	MaxAge   int    `mapstructure:"max_age"`
	BackUp   int    `mapstructure:"max_backups"`
}
