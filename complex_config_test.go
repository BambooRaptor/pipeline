package pipeline_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/BambooRaptor/pipeline"
)

type Config struct {
	AppName string
	Port    int
	APIKey  string
	Timeout time.Duration
}

func NewConfig(appName string, apiKey string) Config {
	return Config{
		AppName: appName,
		Port:    8080,
		APIKey:  apiKey,
		Timeout: time.Second * 10,
	}
}

func SetTimeout(timeout time.Duration) pipeline.Pipe[Config] {
	return func(cfg Config) Config {
		cfg.Timeout = timeout
		return cfg
	}
}

func SetPort(port int) pipeline.Pipe[Config] {
	return func(cfg Config) Config {
		cfg.Port = port
		return cfg
	}
}

func configureApp(cfg Config) Config {
	return cfg
}

func (cfg Config) Matches(ocfg Config) bool {
	return cfg.AppName == ocfg.AppName && cfg.APIKey == ocfg.APIKey && cfg.Timeout == ocfg.Timeout && cfg.Port == ocfg.Port
}

func TestComplexConfig(t *testing.T) {
	config := configureApp(
		pipeline.New(
			SetPort(3000),
			SetTimeout(time.Second*60),
		).Resolve(
			NewConfig("test app", "api-secret"),
		),
	)

	if !config.Matches(Config{"test app", 3000, "api-secret", time.Second * 60}) {
		t.Fatalf("Configs do not match!")
	}
}
