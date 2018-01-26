package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/ans-ashkan/thc/config"
	"github.com/ans-ashkan/thc/twitter"

	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var followingsCmd = &cobra.Command{
	Use:   "followings",
	Short: "get list of following ids",
	Run: func(cmd *cobra.Command, args []string) {
		filepath, _ := cmd.Flags().GetString("output")
		filepath = strings.Replace(filepath, "{date_time}", time.Now().Format("20060102_150405"), 1)

		fmt.Printf("Starting GetFollowings(%s)\n", filepath)

		cfg := config.GetConfig()
		client := twitter.NewClient(cfg.APIKey, cfg.APISecret, cfg.Token, cfg.TokenSecret)
		followingIds, err := client.GetFollowings()
		if err != nil {
			panic(err)
		}

		count := len(followingIds)
		if count > 0 {
			fmt.Printf("Writing %d followings' ids to %s.", count, filepath)
			ids := make([]string, 0, count)
			for _, v := range followingIds {
				ids = append(ids, strconv.FormatInt(v, 10))
			}
			err := ioutil.WriteFile(filepath, []byte(strings.Join(ids, ",")), 0644)
			if err != nil {
				panic(fmt.Errorf("Error writing to file %s. %s", filepath, err))
			}
		} else {
			fmt.Println("No followings to write.")
		}
	},
}

func init() {
	followingsCmd.Flags().StringP("output", "o", "followings_{date_time}.txt", "path to output file")

	RootCmd.AddCommand(followingsCmd)
}
