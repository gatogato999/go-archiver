package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"github.com/DusanKasan/parsemail"
	"github.com/knadh/go-pop3"
)

func processMessages(
	client *pop3.Conn,
	cfg config,
) {
	msgs, err := client.Uidl(0)
	exitOnErr("failed getting msgs", err)
	uploader := newUploader(cfg.CdaEndpoint)
	store, err := loadStore("store.json", cfg.MaxAttempts)
	exitOnErr("failed loading store", err)

	for _, m := range msgs {
		msgBytes, err := client.RetrRaw(m.ID)
		if isErr(fmt.Sprintf("skip: %s, failed retrieve", m.UID), err) {
			continue
		}
		email, err := parsemail.Parse(msgBytes)
		if isErr(fmt.Sprintf("skip: %s, failed paring email", m.UID), err) {
			continue
		}
		if len(email.From) == 0 {
			fmt.Printf("skip: %s, has no sender \n", m.UID)
			continue
		}
		from := email.From[0].Address
		var allAttachsNames []string
		for _, attach := range email.Attachments {
			if attach.Filename == "" {
				continue
			}
			allAttachsNames = append(allAttachsNames, attach.Filename)
			key := Key{UIDL: m.UID, Attachment: attach.Filename}
			if !store.Uploadable(key) {
				log.Printf("skip %s : not Uploadable \n", attach.Filename)
				continue
			}
			data, err := io.ReadAll(attach.Data)
			if isErr(fmt.Sprintf("skip attach %s: reading attach failed", attach.Filename), err) {
				continue
			}
			content := base64.StdEncoding.EncodeToString(data)
			res := uploader.Upload(UploadReq{
				Username: cfg.CdaUser, Password: cfg.CdaPass,
				From: from, Subj: email.Subject,
				Contents: content, Filename: attach.Filename,
				Domain: cfg.CdaDomain,
			})
			log.Printf("uploader response %s -> %+v \n", attach.Filename, res)
			_ = isErr("saving record failed", store.Record(key, res))
		}
		switch {
		case len(allAttachsNames) == 0:
			log.Printf("skip: %s, has no attachments \n", m.UID)
		case store.Resolved(m.UID, allAttachsNames):
			log.Printf("%s, is resolved \n", m.UID)
			if !isErr(fmt.Sprintf("deletion failed, uidl=%s", m.UID), client.Dele(m.ID)) {
				log.Printf("%s, marked for delete \n", m.UID)
			}
		default:
			log.Printf("%s, is not resolved, from: %s \n", m.UID, from)
		}
	}
}
