package conf

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Conf struct {
	Port     string `yaml:"port" env:"PORT" env-default:":1323"`
	DbUrl    string `yaml:"dbUrl" env:"DATABASE_URL" env-required:"true"`
	Server   Server `yaml:"server"`
	JsonLogs bool   `yaml:"jsonLogs" env:"JSON_LOGS" env-default:"false"`
}

type Server struct {
	Url             string `yaml:"url" env:"SERVER_URL" env-required:"true"`
	HealthCheckPath string `yaml:"healthCheckPath" env:"HEALTH_CHECK_PATH" env-default:"/healthz"`
}

var configSingleton Conf
var once sync.Once

func MustGetConfig(configFile ...string) Conf {
	once.Do(func() {
		var err error

		err = cleanenv.ReadEnv(&configSingleton)

		for _, fileName := range configFile {
			err = cleanenv.ReadConfig(fileName, &configSingleton)
		}
		if err != nil {
			panic(err)
		}
	})
	return configSingleton
}
