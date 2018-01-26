package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/ans-ashkan/thc/config"
	"github.com/ans-ashkan/thc/twitter"

	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type filesType []string

func (f filesType) Len() int      { return len(f) }
func (f filesType) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f filesType) Less(i, j int) bool {
	datetimeI, err := time.Parse("followers_20060102_150405.txt", f[i])
	if err != nil {
		panic(fmt.Errorf("error parsing time out of %s, %s", f[i], err))
	}
	datetimeJ, err := time.Parse("followers_20060102_150405.txt", f[j])
	if err != nil {
		panic(fmt.Errorf("error parsing time out of %s, %s", f[j], err))
	}
	return datetimeI.Before(datetimeJ)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "get list of new/un follower ids since the file provided",
	Run: func(cmd *cobra.Command, args []string) {
		lastFilename, _ := cmd.Flags().GetString("last")

		if strings.Index(lastFilename, "{last}") > 0 {
			//find last file in working dir based on followers_{date_time}.txt format
			files, err := filepath.Glob("followers_*.txt")
			if err != nil {
				panic(fmt.Errorf("error finding last followers. %s", err))
			}

			if len(files) == 0 {
				panic(fmt.Errorf("no files found"))
			}
			sort.Sort(sort.Reverse(filesType(files)))
			lastFilename = files[0]
		}

		fileData, err := ioutil.ReadFile(lastFilename)
		if err != nil {
			panic(fmt.Errorf("error reading last followers file %s, %s", lastFilename, err))
		}
		lastFollowerIdsStr := strings.Split(string(fileData), ",")
		lastFollowers := make(map[int64]bool)
		for _, idStr := range lastFollowerIdsStr {
			id, _ := strconv.ParseInt(idStr, 10, 64)
			lastFollowers[id] = true
		}

		newFollowersfilepath, _ := cmd.Flags().GetString("newFollowers_output")
		newFollowersfilepath = strings.Replace(newFollowersfilepath, "{date_time}", time.Now().Format("20060102_150405"), 1)

		unfollowersFilepath, _ := cmd.Flags().GetString("unfollowers_output")
		unfollowersFilepath = strings.Replace(unfollowersFilepath, "{date_time}", time.Now().Format("20060102_150405"), 1)

		fmt.Println("Getting current followers")

		cfg := config.GetConfig()
		client := twitter.NewClient(cfg.APIKey, cfg.APISecret, cfg.Token, cfg.TokenSecret)
		currentFollowerIds, err := client.GetFollowers()
		if err != nil {
			panic(err)
		}

		currentFollowersMap := make(map[int64]bool)
		count := len(currentFollowerIds)
		if count > 0 {
			fmt.Printf("Received %d current followers' ids.\n", count)
			newFollowerIds := make([]string, 0, count)
			for _, v := range currentFollowerIds {
				currentFollowersMap[v] = true
				if _, exists := lastFollowers[v]; !exists {
					newFollowerIds = append(newFollowerIds, strconv.FormatInt(v, 10))
				}
			}
			if len(newFollowerIds) > 0 {
				err := ioutil.WriteFile(newFollowersfilepath, []byte(strings.Join(newFollowerIds, ",")), 0644)
				if err != nil {
					panic(fmt.Errorf("Error writing new followers to file %s. %s", newFollowersfilepath, err))
				}
				fmt.Printf("Written %d new followers to %s\n", len(newFollowerIds), newFollowersfilepath)
			} else {
				fmt.Println("No new followers to write.")
			}

			unfollowerIds := make([]string, 0)
			for lastFollowerID := range lastFollowers {
				if _, exists := currentFollowersMap[lastFollowerID]; !exists {
					unfollowerIds = append(unfollowerIds, strconv.FormatInt(lastFollowerID, 10))
				}
			}

			if len(unfollowerIds) > 0 {
				err := ioutil.WriteFile(unfollowersFilepath, []byte(strings.Join(unfollowerIds, ",")), 0644)
				if err != nil {
					panic(fmt.Errorf("Error writing unfollowers to file %s. %s", unfollowersFilepath, err))
				}
				fmt.Printf("Written %d unfollowers to %s\n", len(unfollowerIds), unfollowersFilepath)
			} else {
				fmt.Println("No unfollowers to write.")
			}
		} else {
			fmt.Println("No unfollowers to write.")
		}
	},
}

func init() {
	diffCmd.Flags().StringP("last", "f", "followers_{last}.txt", "path to last followers file")
	diffCmd.Flags().StringP("newFollowers_output", "n", "newFollowers_{date_time}.txt", "path to new followers output file")
	diffCmd.Flags().StringP("unfollowers_output", "u", "unfollowers_{date_time}.txt", "path to unfollowers output file")

	RootCmd.AddCommand(diffCmd)
}
