package config

type ProducerConf struct {
	URI   string `mapstructure:"uri"`
	Queue string `mapstructure:"queue"`
}
