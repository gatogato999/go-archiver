package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/knadh/go-pop3"
)

func exitOnErr(details string, err error) {
	if isErr(details, err) { log.Fatal("stoping execution") }
}

func isErr(details string, err error) bool {
	if err != nil {
		log.Printf("%s : %v \n", details, err)
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

func loadConfig() (cfg config) {
	exitOnErr("failed to read .evn", cleanenv.ReadConfig(".env", &cfg))
	return cfg
}

func connectClient(cfg config) (client *pop3.Conn ) {
	client, err := pop3.New(pop3.Opt{ Host: cfg.PopHost, Port: 995, TLSEnabled: true, }).NewConn()
	exitOnErr("client creation failed", err)
	exitOnErr("client auth failed", client.Auth(cfg.PopUser, cfg.PopPass))
	return client 
}

func initLog()  *os.File {
os.MkdirAll("logs", 0755)
f, err := os.OpenFile( "logs/"+time.Now().Format("02-01-2006")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644,)
exitOnErr("can't open/create log folder/file", err)
log.SetOutput(io.MultiWriter(os.Stdout, f))
return f
}
