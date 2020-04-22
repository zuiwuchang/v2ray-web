package cmd

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/spf13/cobra"
	"gitlab.com/king011/v2ray-web/db/data"
	"gitlab.com/king011/v2ray-web/utils"
)

func init() {
	var filename string
	basePath := utils.BasePath()
	var name, password string
	cmd := &cobra.Command{
		Use:   `password`,
		Short: `change web user password`,
		Run: func(cmd *cobra.Command, args []string) {
			db, e := bolt.Open(filename, 0600, nil)
			if e != nil {
				log.Fatalln(e)
			}
			if name == "" {
				log.Fatalln("name not support empty")

			}
			if password == "" {
				log.Fatalln("password not support empty")
			}

			pwd := sha512.Sum512([]byte(password))
			e = db.Update(func(t *bolt.Tx) (e error) {
				bucket := t.Bucket(data.UserBucket)
				if bucket == nil {
					e = fmt.Errorf("bucket not exist : %s", data.UserBucket)
					return
				}
				key := utils.StringToBytes(name)
				e = bucket.Put(key, utils.StringToBytes(hex.EncodeToString(pwd[:])))
				return
			})
			if e != nil {
				log.Fatalln(e)
			}
			fmt.Println("password reset complete")
		},
	}
	flasg := cmd.Flags()
	flasg.StringVarP(&filename,
		"filename", "f",
		utils.Abs(basePath, "v2ray-web.db"),
		"database source filename",
	)
	flasg.StringVarP(&name,
		"name", "n",
		"",
		"user name",
	)
	flasg.StringVarP(&password,
		"password", "p",
		"",
		"user password",
	)
	rootCmd.AddCommand(cmd)
}
