package main

import (
	"log"

	_ "github.com/v2fly/v2ray-core/v4/main/distro/all"
	_ "gitlab.com/king011/v2ray-web/assets/static/statik"
	_ "gitlab.com/king011/v2ray-web/assets/zh-Hans/statik"
	_ "gitlab.com/king011/v2ray-web/assets/zh-Hant/statik"
	"gitlab.com/king011/v2ray-web/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}
