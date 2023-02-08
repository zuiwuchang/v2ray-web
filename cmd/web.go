package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/king011/v2ray-web/cmd/web"
	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/cookie"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/single/upgrade"
	"gitlab.com/king011/v2ray-web/utils"
)

func init() {
	var filename, addr string
	var noupgrade, debug, none bool
	basePath := utils.BasePath()
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
			e = cnf.Format(basePath)
			if e != nil {
				log.Fatalln(e)
			}
			if addr != `` {
				cnf.HTTP.Addr = addr
			}

			// init logger
			e = logger.Init(basePath, &cnf.Logger)
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
			if !noupgrade {
				go upgrade.DefaultUpgrade().Serve()
			}
			web.Run(cnf, debug, none)
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename,
		"config", "c",
		utils.Abs(basePath, "v2ray-web.jsonnet"),
		"Config file for Web",
	)
	flags.BoolVarP(&debug,
		"debug", "d",
		false,
		"Run as debug",
	)
	flags.StringVarP(&addr,
		"addr", "a",
		"",
		"Web listen address",
	)
	flags.BoolVar(&noupgrade, `no-upgrade`,
		false,
		`Disable automatic upgrades`,
	)
	flags.BoolVar(&none, `none`, false, `Do not start v2ray nor set iptables`)
	rootCmd.AddCommand(cmd)
}
