package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xtls/xray-core/core"
	"gitlab.com/king011/v2ray-web/single/upgrade"
)

func init() {
	var filename string
	var test, noupgrade bool
	cmd := &cobra.Command{
		Use:   `start`,
		Short: `Start V2Ray server.`,
		Run: func(cmd *cobra.Command, args []string) {
			f, e := os.Open(filename)
			if e != nil {
				log.Fatalln(e)
			}
			cnf, e := core.LoadConfig(`json`, f)
			f.Close()
			if e != nil {
				log.Fatalln(e)
			}
			server, e := core.New(cnf)
			if e != nil {
				log.Fatalln(e)
			}
			if test {
				return
			}
			e = server.Start()
			if e != nil {
				log.Fatalln(e)
			}
			defer server.Close()
			if !noupgrade {
				go upgrade.DefaultUpgrade().Serve()
			}

			{
				osSignals := make(chan os.Signal, 1)
				signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
				<-osSignals
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVarP(&filename,
		"config", "c",
		"",
		"Config file for V2Ray",
	)
	flags.BoolVarP(&test,
		"test", "t",
		false,
		"Test config file only, without launching V2Ray server.",
	)
	flags.BoolVar(&noupgrade, `no-upgrade`,
		false,
		`disable automatic upgrades`,
	)
	rootCmd.AddCommand(cmd)
}
