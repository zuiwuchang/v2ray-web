package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/king011/v2ray-web/cmd/web"
	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/utils"
)

func init() {
	var filename string
	basePaht := utils.BasePath()
	cmd := &cobra.Command{
		Use:   `web`,
		Short: `Start V2Ray web control server.`,
		Run: func(cmd *cobra.Command, args []string) {
			// load configure
			cnf := configure.Single()
			e := cnf.Load(filename)
			if e != nil {
				log.Fatalln(e)
			}
			e = cnf.Format(basePaht)
			if e != nil {
				log.Fatalln(e)
			}

			// init logger
			e = logger.Init(basePaht, &cnf.Logger)
			if e != nil {
				log.Fatalln(e)
			}

			// init cookie
			e = cookie.Init(cnf.HTTP.Secret, cnf.HTTP.Maxage)
			if e != nil {
				log.Fatalln(e)
			}

			// init db
			e = manipulator.Init(&cnf.Database)
			if e != nil {
				log.Fatalln(e)
			}

			web.Run(cnf)
		},
	}
	flasg := cmd.Flags()
	flasg.StringVarP(&filename,
		"config", "c",
		utils.Abs(basePaht, "v2ray-web.jsonnet"),
		"Config file for Web",
	)
	rootCmd.AddCommand(cmd)
}
