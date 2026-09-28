package main

import (
	"encoding/base64"
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
	exitOnErr("failed getting msgs", err)
	uploader := newUploader(cfg.CdaEndpoint)
	store, err := loadStore("store.json", 4)
	exitOnErr("failed loading store", err)
	for _, m := range msgs {
		log.Printf(
			"\n---------------- %s -----------------\n",
			m.UID,
		)
		msgBytes, err := client.RetrRaw(m.ID)
		if isErr(fmt.Sprintf("skip: %s, failed retrieve", m.UID), err) {
			continue
		}
		email, err := parsemail.Parse(msgBytes)
		if isErr(fmt.Sprintf("skip: %s, failed paring email", m.UID), err) {
			continue
		}
		from := "unknown"
		for _, f := range email.From {
			from = f.Address
		}
		log.Printf("From %s ", from)
		log.Printf("subj %s ", email.Subject)

		if len(email.Attachments) == 0 {
			fmt.Printf("skip: %s, has not attachments", m.UID)
			continue
		}
		var allAttachsNames []string
		for _, attach := range email.Attachments {
			if attach.Filename == "" {
				continue
			}

			log.Printf("\n---------------- %s - %s -----------------\n", m.UID, attach.Filename)
			allAttachsNames = append(allAttachsNames, attach.Filename)
			key := Key{
				UIDL:       m.UID,
				Attachment: attach.Filename,
			}
			if !store.Uploadable(key) {
				fmt.Printf("skip %s %s : not Uploadable", m.UID, attach.Filename)
				continue
			}
			data, err := io.ReadAll(attach.Data)
			if isErr(fmt.Sprintf("skip attach %s: reading attach failed", attach.Filename), err) {
				continue
			}
			if len(data) > 1 {
				content := base64.StdEncoding.EncodeToString(data)
				res := uploader.Upload(UploadReq{
					Username: cfg.CdaUser, Password: cfg.CdaPass,
					From: from, Subj: email.Subject,
					Contents: content, Filename: attach.Filename,
					Domain: cfg.CdaDomain,
				})
				err = store.Record(key, res)
				isErr("saving record failed", err)
				log.Printf(
					"\n uidl: %s, attach : %s, uploader response : %s",
					m.UID,
					attach.Filename,
					res,
				)
			}
		}
		switch {
		case len(allAttachsNames) == 0:
			fmt.Printf("\n skip: %s, has not attachments", m.UID)
		case store.Resolved(m.UID, allAttachsNames):
			fmt.Printf("\n %s, is resolved", m.UID)
			err = client.Dele(m.ID)
			if !isErr(fmt.Sprintf("deletion failed, uidl=%s", m.UID), err) {
				fmt.Printf("\n %s, marked for delete", m.UID)
			}
		default:
			fmt.Printf("\n %s, is not resolved, from: %s", m.UID, from)
		}
	}
}
