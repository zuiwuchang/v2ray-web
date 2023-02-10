package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/utils"
)

func init() {
	var (
		db string
		settings, strategy, v2ray, subscription,
		iptables,
		iptablesView, iptablesClear, iptablesInit,
		users, proxy, last string
		basePath = utils.BasePath()
	)
	cmd := &cobra.Command{
		Use:   `import`,
		Short: `Import settings`,
		Run: func(cmd *cobra.Command, args []string) {
			if settings != `` {
				var v data.Settings
				importJSON(settings, &v)

				initDB(basePath, db)
				var mSettings manipulator.Settings
				e := mSettings.Put(&v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - settings from:`, settings)
			}
			if strategy != `` {
				var v []*data.Strategy
				importJSON(strategy, &v)

				initDB(basePath, db)
				var mStrategy manipulator.Strategy
				e := mStrategy.Import(v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - strategy from:`, strategy)
			}

			if v2ray != `` {
				v := importText(v2ray)

				initDB(basePath, db)
				var mSettings manipulator.Settings
				e := mSettings.PutV2ray(v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - v2ray from:`, v2ray)
			}

			if subscription != `` {
				var v []*data.Subscription
				importJSON(subscription, &v)

				initDB(basePath, db)
				var mSubscription manipulator.Subscription
				e := mSubscription.Import(v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - subscription from:`, subscription)
			}

			if iptables != `` || iptablesView != `` || iptablesClear != `` || iptablesInit != `` {
				initDB(basePath, db)
				var (
					mSettings manipulator.Settings
					v         *data.IPTables
					e         error
				)
				if iptables == `` {
					v, e = mSettings.GetIPtables()
					if e != nil {
						log.Fatalln(e)
					}
				} else {
					v = &data.IPTables{}
					importJSON(iptables, v)
				}

				if iptablesView != `` {
					v.View = importText(iptablesView)
				}
				if iptablesClear != `` {
					v.Clear = importText(iptablesClear)
				}
				if iptablesInit != `` {
					v.Init = importText(iptablesInit)
				}

				e = mSettings.PutIPtables(v)
				if e != nil {
					log.Fatalln(e)
				}
				if iptables != `` {
					fmt.Println(` - iptables from:`, iptables)
				}
				if iptablesView != `` {
					fmt.Println(` - iptables view from:`, iptablesView)
				}
				if iptablesClear != `` {
					fmt.Println(` - iptables clear from:`, iptablesClear)
				}
				if iptablesInit != `` {
					fmt.Println(` - iptables init from:`, iptablesInit)
				}
			}
			if users != `` {
				var v []data.UserRaw
				importJSON(users, &v)

				initDB(basePath, db)
				var mUser manipulator.User
				e := mUser.Import(v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - users from:`, users)
			}
			if proxy != `` {
				var v []*data.Element
				importJSON(proxy, &v)

				initDB(basePath, db)
				var mElement manipulator.Element

				e := mElement.Import(v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - proxy from:`, proxy)
			}
			if last != `` {
				var v data.Element
				importJSON(last, &v)

				initDB(basePath, db)
				var mSettings manipulator.Settings

				e := mSettings.PutLast(&v)
				if e != nil {
					log.Fatalln(e)
				}
				fmt.Println(` - last from:`, last)
			}
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&db, `db`,
		utils.Abs(basePath, "v2ray-web.db"),
		`db path`,
	)
	flags.StringVarP(&settings, `settings`,
		`s`,
		``,
		`if non-empty import settings from this file`,
	)
	flags.StringVar(&strategy, `strategy`,
		``,
		`if non-empty import strategy from this file`,
	)
	flags.StringVarP(&v2ray, `v2ray`,
		`v`,
		``,
		`if non-empty import v2ray configure template from this file`,
	)
	flags.StringVar(&subscription, `subscription`,
		``,
		`if non-empty import subscription from this file`,
	)

	flags.StringVar(&iptables, `iptables`,
		``,
		`if non-empty import v2ray iptables template from this file`,
	)
	flags.StringVar(&iptablesView, `iptables-view`,
		``,
		`if non-empty import v2ray iptables view template from this file`,
	)
	flags.StringVar(&iptablesClear, `iptables-clear`,
		``,
		`if non-empty import v2ray iptables clear template from this file`,
	)
	flags.StringVar(&iptablesInit, `iptables-init`,
		``,
		`if non-empty import v2ray iptables init template from this file`,
	)
	flags.StringVarP(&users, `users`,
		`u`,
		``,
		`if non-empty import users from this file`,
	)
	flags.StringVarP(&proxy, `proxy`,
		`p`,
		``,
		`if non-empty import proxy element from this file`,
	)
	flags.StringVarP(&last, `last`,
		`l`,
		``,
		`if non-empty import last start proxy`,
	)
	rootCmd.AddCommand(cmd)
}
func importJSON(filepath string, v interface{}) {
	b, e := os.ReadFile(filepath)
	if e != nil {
		log.Fatalln(e)
	}
	e = json.Unmarshal(b, v)
	if e != nil {
		log.Fatalln(e)
	}
}
func importText(filepath string) string {
	b, e := os.ReadFile(filepath)
	if e != nil {
		log.Fatalln(e)
	}
	return string(b)
}
