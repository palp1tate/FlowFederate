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
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Database   int    `json:"database"`
	Expiration int    `json:"expiration"`
}

type AliSmsConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

type QiNiuYunConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
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
