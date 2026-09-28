package main

import (
	"log"

	"github.com/knadh/go-pop3"
)

func main() {
	cfg := loadConfig()
	client, err := pop3.New(pop3.Opt{
		Host:       cfg.PopHost,
		Port:       995,
		TLSEnabled: true,
	}).NewConn()
	exitOnErr("client creation failed", err)
	defer closeConn(*client)
	err = client.Auth(cfg.PopUser, cfg.PopPass)
	exitOnErr("client auth failed", err)
	msgs, err := client.Uidl(0)
	log.Println(msgs)
}
