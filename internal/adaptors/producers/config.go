package producers

import (
	"fmt"
	"golang-skeleton/pkg/env"
	"golang-skeleton/pkg/errors"
	"log"
)

type Config struct {
	BootStrapServers []string `env:"PRODUCER_BOOTSTRAP_SERVERS"`
	ClientID         string   `env:"PRODUCER_CLIENT_ID" envDefault:"golang-restful-streaming-producer"`
}

func (k *Config) Register() error {
	err := env.Parse(k)
	if err != nil {
		return fmt.Errorf("error loading abstractions Config, err: %v", err)
	}

	return nil
}

func (k *Config) Validate() error {
	if k.BootStrapServers == nil {
		return errors.New(`kafka brokers can not be empty`)
	}
	if k.ClientID == `` {
		return errors.New(`producer client id can not be empty`)
	}

	return nil
}

func (k *Config) Print() interface{} {
	defer log.Println("k-abstractions configs loaded")

	return k
}
