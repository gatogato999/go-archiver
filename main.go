package main

func main() {
	f := initLog()
	defer f.Close()
	cfg := loadConfig()
	client := connectClient(cfg)
	defer func() { exitOnErr("client quit failed, msg deletion may have failed", client.Quit()) }()
	processMessages(client, cfg)
}
