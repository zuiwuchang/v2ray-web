package main

import (
	"log"

	_ "github.com/v2fly/v2ray-core/v4/main/distro/all"
	"gitlab.com/king011/v2ray-web/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
