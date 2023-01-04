package main

import (
	"log"

	_ "github.com/xtls/xray-core/main/distro/all"

	"gitlab.com/king011/v2ray-web/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
