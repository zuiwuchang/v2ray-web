package main

import (
	"log"

	"gitlab.com/king011/v2ray-web/cmd"
	_ "v2ray.com/core/main/distro/all"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
