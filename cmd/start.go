package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"v2ray.com/core"
)

func init() {
	var filename string
	var test bool
	cmd := &cobra.Command{
		Use:   `start`,
		Short: `Start V2Ray server.`,
		Run: func(cmd *cobra.Command, args []string) {
			f, e := os.Open(filename)
			cnf, e := core.LoadConfig(`json`, filename, f)
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
			{
				osSignals := make(chan os.Signal, 1)
				signal.Notify(osSignals, os.Interrupt, os.Kill, syscall.SIGTERM)
				<-osSignals
			}
		},
	}
	flasg := cmd.Flags()
	flasg.StringVarP(&filename,
		"config", "c",
		"",
		"Config file for V2Ray",
	)
	flasg.BoolVarP(&test,
		"test", "t",
		false,
		"Test config file only, without launching V2Ray server.",
	)
	rootCmd.AddCommand(cmd)
}
