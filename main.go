package main

import (
	"fmt"
	"io"
	"log"

	"github.com/DusanKasan/parsemail"
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
	for _, m := range msgs {
		msgBytes, err := client.RetrRaw(m.ID)
		if isErr(fmt.Sprintf("skip: %d, failed retrieve", m.UID), err) {
			continue
		}
		email, err := parsemail.Parse(msgBytes)
		if isErr(fmt.Sprintf("skip: %d, failed paring email", m.UID), err) {
			continue
		}
		log.Printf("From %s ", email.From)
		log.Printf("subj %s ", email.Subject)
		for _, attach := range email.Attachments {
			log.Printf("file name %s ", attach.Filename)
			log.Printf("file ContentType %s ", attach.ContentType)
			data, err := io.ReadAll(attach.Data)
			if isErr("skip %d: reading attach failed ", err) {
				continue
			}
			if len(data) > 20 {
				log.Printf("file data %s ", string(data[:20]))
			}
		}

	}
}
