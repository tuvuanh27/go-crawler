package config

import (
	"flag"
	"fmt"
	"github.com/tuvuanh27/go-crawler/internal/pkg/grpc"
	echoserver "github.com/tuvuanh27/go-crawler/internal/pkg/http/echo/server"
	mongodriver "github.com/tuvuanh27/go-crawler/internal/pkg/mongo-driver"
	"github.com/tuvuanh27/go-crawler/internal/pkg/otel"
	"github.com/tuvuanh27/go-crawler/internal/pkg/rabbitmq"
	"github.com/tuvuanh27/go-crawler/internal/services/worker/service"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "products write microservice config path")
}

type Config struct {
	ServiceName    string                        `mapstructure:"serviceName"`
	Rabbitmq       *rabbitmq.RabbitMQConfig      `mapstructure:"rabbitmq"`
	Echo           *echoserver.EchoConfig        `mapstructure:"echo"`
	Grpc           *grpc.GrpcConfig              `mapstructure:"grpc"`
	MongoConfig    *mongodriver.MongoConfig      `mapstructure:"mongoConfig"`
	Jaeger         *otel.JaegerConfig            `mapstructure:"jaeger"`
	OpenAliexpress *service.OpenAliexpressConfig `mapstructure:"openAliexpress"`
}

type Context struct {
	Timeout int `mapstructure:"timeout"`
}

func InitConfig() (*Config, *otel.JaegerConfig, *mongodriver.MongoConfig,
	*grpc.GrpcConfig, *echoserver.EchoConfig, *rabbitmq.RabbitMQConfig, *service.OpenAliexpressConfig, error) {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	if configPath == "" {
		configPathFromEnv := os.Getenv("CONFIG_PATH")
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			//https://stackoverflow.com/questions/31873396/is-it-possible-to-get-the-current-root-of-package-structure-as-a-string-in-golan
			//https://stackoverflow.com/questions/18537257/how-to-get-the-directory-of-the-currently-running-file
			d, err := dirname()
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, err
			}

			configPath = d
		}
	}

	cfg := &Config{}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.AddConfigPath(configPath)
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, nil, nil, nil, nil, nil, nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return cfg, cfg.Jaeger, cfg.MongoConfig, cfg.Grpc, cfg.Echo, cfg.Rabbitmq, cfg.OpenAliexpress, nil
}

func GetMicroserviceName(serviceName string) string {
	return fmt.Sprintf("%s", strings.ToUpper(serviceName))
}

func filename() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}
	return filename, nil
}

func dirname() (string, error) {
	filename, err := filename()
	if err != nil {
		return "", err
	}
	return filepath.Dir(filename), nil
}
