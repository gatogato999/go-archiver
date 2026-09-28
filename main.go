package main

import (
	"github.com/knadh/go-pop3"
)

func main() {
	cfg := loadConfig()
	client, err := pop3.New(pop3.Opt{
		Host: cfg.PopHost, Port: 995, TLSEnabled: true,
	}).NewConn()
	exitOnErr("client creation failed", err)
	defer func() { exitOnErr("client quit failed, msg deletion may have failed", client.Quit()) }()
	exitOnErr("client auth failed", client.Auth(cfg.PopUser, cfg.PopPass))
	msgs, err := client.Uidl(0)
	exitOnErr("failed getting msgs", err)
	uploader := newUploader(cfg.CdaEndpoint)
	store, err := loadStore("store.json", cfg.MaxAttempts)
	exitOnErr("failed loading store", err)
	for _, m := range msgs {
		processMessages(client, store, uploader, cfg, m)
	}
}
