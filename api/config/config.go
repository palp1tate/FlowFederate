package config

type ApiConfig struct {
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags"`
}
type ServiceConfig struct {
	Edge string `yaml:"edge"`
}

type JWTConfig struct {
	SigningKey string `yaml:"signingKey"`
	Expiration int    `yaml:"expiration"`
}

type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Database   int    `yaml:"database"`
	Expiration int    `yaml:"expiration"`
}

type AliSmsConfig struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	SignName        string `yaml:"signName"`
	TemplateCode    string `yaml:"templateCode"`
}

type QiNiuYunConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ServerConfig struct {
	Api        ApiConfig      `yaml:"api"`
	Service    ServiceConfig  `yaml:"service"`
	MySQL      MySQLConfig    `yaml:"mysql"`
	JWT        JWTConfig      `yaml:"jwt"`
	Consul     ConsulConfig   `yaml:"consul"`
	Redis      RedisConfig    `yaml:"redis"`
	AliSms     AliSmsConfig   `yaml:"sms"`
	QiNiuYun   QiNiuYunConfig `yaml:"qiNiuYun"`
	CommonName string         `yaml:"commonName"`
	Debug      bool           `yaml:"debug"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"userService"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataId"`
	Group     string `mapstructure:"group"`
}
