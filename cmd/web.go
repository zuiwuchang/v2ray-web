package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	var filename string
	cmd := &cobra.Command{
		Use:   `web`,
		Short: `Start V2Ray web control server.`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	flasg := cmd.Flags()
	flasg.StringVarP(&filename,
		"config", "c",
		"",
		"Config file for V2Ray",
	)
	rootCmd.AddCommand(cmd)
}
