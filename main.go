package main

func main() {
	cfg := loadConfig()
	client, err := connectClient(cfg)
	exitOnErr("failed connecting", err)
	defer func() { exitOnErr("client quit failed, msg deletion may have failed", client.Quit()) }()
	processMessages(client, cfg)
}
