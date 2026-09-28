package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type AttachStatus string

const (
	Uploaded        AttachStatus = "Uploaded"
	SkippedTooSmall AttachStatus = "SkippedTooSmall"
	FailedRetryable AttachStatus = "FailedRetryable"
	FailedPermanent AttachStatus = "FailedPermanent"
)

type UploadReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"email"`
	Subj     string `json:"title"`
	Contents string `json:"contents"`
	Filename string `json:"filename"`
	Domain   string `json:"domain"`
}

type uploadRes struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Uploader struct {
	client *http.Client
	url    string
}

func newUploader(url string) *Uploader {
	return &Uploader{client: &http.Client{Timeout: 30 * time.Second}, url: url}
}

func (u *Uploader) Upload(req UploadReq) AttachStatus {
	body, err := json.Marshal(req)
	if isErr("marshal failed", err) {
		return FailedRetryable
	}
	resp, err := u.client.Post(u.url, "application/json", bytes.NewReader(body))
	if isErr("call to uploader failed", err) {
		return FailedRetryable
	}
	defer resp.Body.Close()
	var r uploadRes
	err = json.NewDecoder(resp.Body).Decode(&r)
	if isErr("Upload: deserialization failed", err) {
		return FailedRetryable
	}
	switch {
	case r.Message == "Empty data" || r.Message == "Empty attachment" || r.Message == "Ignoring small attachments":
		return SkippedTooSmall
	case r.Success:
		return Uploaded
	}
	return FailedRetryable
}
