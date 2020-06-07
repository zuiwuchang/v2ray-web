package main

import (
	"log"

	_ "gitlab.com/king011/v2ray-web/assets/zh-Hans/statik"
	_ "gitlab.com/king011/v2ray-web/assets/zh-Hant/statik"
	"gitlab.com/king011/v2ray-web/cmd"
	_ "v2ray.com/core/main/distro/all"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
