package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/single/upgrade"
	"gitlab.com/king011/v2ray-web/version"
)

func init() {
	var (
		yes   bool
		loglv string
	)

	cmd := &cobra.Command{
		Use:   `upgrade`,
		Short: `Upgrade to the latest version`,
		Run: func(cmd *cobra.Command, args []string) {
			logger.InitConsole(strings.ToLower(strings.TrimSpace(loglv)))
			upgraded, ver, e := upgrade.DefaultUpgrade().Do(yes)
			if e == nil {
				if upgraded {
					fmt.Println(`upgrade success:`, version.Version, `->`, ver)
				} else {
					fmt.Println(`already the latest version:`, version.Version)
				}
			} else {
				log.Fatalln(e)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&loglv, `log`,
		`info`,
		`log level [debug info warn error dpanic panic fatal]`,
	)
	flags.BoolVarP(&yes, `yes`,
		`y`,
		false,
		`automatic yes to prompts`,
	)
	rootCmd.AddCommand(cmd)
}
