package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/knadh/go-pop3"
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

func closeConn(c pop3.Conn) {
	err := c.Quit()
	exitOnErr("client quit failed, delation may have failed", err)
	log.Println("quitting")
}

type config struct {
	PopHost     string `env:"POP_DOMAIN"   env-required:"true"`
	PopUser     string `env:"POP_USER"     env-required:"true"`
	PopPass     string `env:"POP_PASS"     env-required:"true"`
	CdaUser     string `env:"CDA_USER"     env-required:"true"`
	CdaPass     string `env:"CDA_PASS"     env-required:"true"`
	CdaDomain   string `env:"CDA_DOMAIN"                       env-default:"golang"`
	CdaEndpoint string `env:"CDA_ENDPOINT"                     env-default:"http://localhost:10032/cda/UploadAttachment"`
}

func loadConfig() config {
	var cfg config
	err := cleanenv.ReadConfig(".env", &cfg)
	exitOnErr("failed to read .evn", err)
	return cfg
}
