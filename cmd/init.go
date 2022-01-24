package cmd

import (
	"fmt"
	"runtime"

	core "github.com/v2fly/v2ray-core/v4"
	"gitlab.com/king011/v2ray-web/version"

	"github.com/spf13/cobra"
)

const (
	// App .
	App = "v2ray-web"
)

var v bool
var rootCmd = &cobra.Command{
	Use:   App,
	Short: "generate tools create full",
	Run: func(cmd *cobra.Command, args []string) {
		if v {
			fmt.Println(runtime.GOOS, runtime.GOARCH, runtime.Version())
			fmt.Println(version.Version)
			fmt.Println(version.Commit)
			fmt.Println(version.Date)
			version := core.VersionStatement()
			for _, s := range version {
				fmt.Println(s)
			}
		} else {
			fmt.Println(App)
			fmt.Println(runtime.GOOS, runtime.GOARCH, runtime.Version())
			fmt.Println(version.Version)
			fmt.Println(version.Commit)
			fmt.Println(version.Date)
			version := core.VersionStatement()
			for _, s := range version {
				fmt.Println(s)
			}
			fmt.Printf(`Use "%v --help" for more information about this program.
`, App)
		}
	},
}

func init() {
	flags := rootCmd.Flags()
	flags.BoolVarP(&v,
		"version",
		"v",
		false,
		"show version",
	)
}

// Execute run command
func Execute() error {
	return rootCmd.Execute()
}
