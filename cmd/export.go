package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	klogger "gitlab.com/king011/king-go/log/logger.zap"
	"gitlab.com/king011/v2ray-web/configure"
	"gitlab.com/king011/v2ray-web/db/manipulator"
	"gitlab.com/king011/v2ray-web/logger"
	"gitlab.com/king011/v2ray-web/utils"
)

var once sync.Once

func initDB(basePath, path string) {
	once.Do(func() {
		e := logger.Init(basePath, &klogger.Options{
			Level: "error",
		})
		if e != nil {
			log.Fatalln(e)
		}
		e = manipulator.Init(
			&configure.Database{
				Source: path,
			},
			&bolt.Options{
				Timeout: time.Second * 5,
			},
		)
		if e != nil {
			log.Fatalln(e)
		}
	})
}
func init() {
	var (
		db string
		settings, strategy, v2ray, subscription,
		iptables,
		iptablesView, iptablesClear, iptablesInit,
		users, proxy string
		basePath = utils.BasePath()
	)

	cmd := &cobra.Command{
		Use:   `export`,
		Short: `Export settings`,
		Run: func(cmd *cobra.Command, args []string) {
			if settings != `` {
				initDB(basePath, db)
				var mSettings manipulator.Settings
				v, e := mSettings.Get()
				if e != nil {
					log.Fatalln(e)
				}
				exportJSON(settings, v)
				fmt.Println(` - settings to:`, settings)
			}
			if strategy != `` {
				initDB(basePath, db)
				var mStrategy manipulator.Strategy
				v, e := mStrategy.List()
				if e != nil {
					log.Fatalln(e)
				}
				exportJSON(strategy, v)
				fmt.Println(` - strategy to:`, strategy)
			}
			if v2ray != `` {
				initDB(basePath, db)
				var mSettings manipulator.Settings
				v, e := mSettings.GetV2ray()
				if e != nil {
					log.Fatalln(e)
				}
				exportText(v2ray, v)
				fmt.Println(` - v2ray to:`, v2ray)
			}
			if subscription != `` {
				initDB(basePath, db)
				var mSubscription manipulator.Subscription
				v, e := mSubscription.List()
				if e != nil {
					log.Fatalln(e)
				}
				exportJSON(subscription, v)
				fmt.Println(` - subscription to:`, subscription)
			}
			if iptables != `` || iptablesView != `` || iptablesClear != `` || iptablesInit != `` {
				initDB(basePath, db)
				var mSettings manipulator.Settings
				v, e := mSettings.GetIPtables()
				if e != nil {
					log.Fatalln(e)
				}
				if iptables != `` {
					exportJSON(iptables, v)
					fmt.Println(` - iptables to:`, iptables)
				}
				if iptablesView != `` {
					exportText(iptablesView, v.View)
					fmt.Println(` - iptables view to:`, iptablesView)
				}
				if iptablesClear != `` {
					exportText(iptablesClear, v.Clear)
					fmt.Println(` - iptables clear to:`, iptablesClear)
				}
				if iptablesInit != `` {
					exportText(iptablesInit, v.Init)
					fmt.Println(` - iptables init to:`, iptablesInit)
				}
			}
			if users != `` {
				initDB(basePath, db)
				var mUser manipulator.User
				v, e := mUser.ListRaw()
				if e != nil {
					log.Fatalln(e)
				}
				exportJSON(users, v)
				fmt.Println(` - users to:`, users)
			}
			if proxy != `` {
				initDB(basePath, db)
				var mElement manipulator.Element
				v, _, e := mElement.List()
				if e != nil {
					log.Fatalln(e)
				}
				exportJSON(proxy, v)
				fmt.Println(` - proxy to:`, proxy)
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
		`if non-empty export settings to this file`,
	)
	flags.StringVar(&strategy, `strategy`,
		``,
		`if non-empty export strategy to this file`,
	)
	flags.StringVarP(&v2ray, `v2ray`,
		`v`,
		``,
		`if non-empty export v2ray configure template to this file`,
	)
	flags.StringVar(&subscription, `subscription`,
		``,
		`if non-empty export subscription to this file`,
	)

	flags.StringVar(&iptables, `iptables`,
		``,
		`if non-empty export v2ray iptables template to this file`,
	)
	flags.StringVar(&iptablesView, `iptables-view`,
		``,
		`if non-empty export v2ray iptables view template to this file`,
	)
	flags.StringVar(&iptablesClear, `iptables-clear`,
		``,
		`if non-empty export v2ray iptables clear template to this file`,
	)
	flags.StringVar(&iptablesInit, `iptables-init`,
		``,
		`if non-empty export v2ray iptables init template to this file`,
	)
	flags.StringVarP(&users, `users`,
		`u`,
		``,
		`if non-empty export users to this file`,
	)
	flags.StringVarP(&proxy, `proxy`,
		`p`,
		``,
		`if non-empty export proxy element to this file`,
	)

	rootCmd.AddCommand(cmd)
}
func exportJSON(filepath string, v interface{}) {
	b, e := json.MarshalIndent(v, "", "\t")
	if e != nil {
		log.Fatalln(e)
	}

	e = os.WriteFile(filepath, b, 0666)
	if e != nil {
		log.Fatalln(e)
	}
}
func exportText(filepath, text string) {
	e := os.WriteFile(filepath, utils.StringToBytes(text), 0666)
	if e != nil {
		log.Fatalln(e)
	}
}
