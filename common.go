package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

func exitOnErr(details string, err error) {
	if isErr(details, err) {
		log.Fatal("")
	}
}

func isErr(details string, err error) bool {
	if err != nil {
		log.Printf("%s : %v", details, err)
		return true
	}
	return false
}

type config struct {
	PopHost     string `env:"POP_DOMAIN"   env-required:"true"`
	PopUser     string `env:"POP_USER"     env-required:"true"`
	PopPass     string `env:"POP_PASS"     env-required:"true"`
	CdaUser     string `env:"CDA_USER"     env-required:"true"`
	CdaPass     string `env:"CDA_PASS"     env-required:"true"`
	CdaDomain   string `env:"CDA_DOMAIN"                       env-default:"golang"`
	CdaEndpoint string `env:"CDA_ENDPOINT"                     env-default:"http://localhost:10032/cda/UploadAttachment"`
	MaxAttempts int    `env:"MAX_ATTEMPTS"                     env-default:"5"`
}

func loadConfig() config {
	var cfg config
	exitOnErr("failed to read .evn", cleanenv.ReadConfig(".env", &cfg))
	return cfg
}
