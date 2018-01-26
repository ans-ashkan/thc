package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ans-ashkan/thc/config"
	"github.com/ans-ashkan/thc/twitter"

	"strconv"
	"strings"
	"time"

	"path/filepath"

	"encoding/json"

	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "get list of users from ids in provided file",
	Run: func(cmd *cobra.Command, args []string) {
		inputPath, _ := cmd.Flags().GetString("input")
		if inputPath == "" {
			panic(fmt.Errorf("input path is not set"))
		}
		basePath := filepath.Base(inputPath)
		inputFileName := strings.TrimSuffix(basePath, filepath.Ext(basePath))

		outputPath, _ := cmd.Flags().GetString("output")
		outputPath = strings.Replace(outputPath, "{input_file_name}", inputFileName, 1)
		outputPath = strings.Replace(outputPath, "{date_time}", time.Now().Format("20060102_150405"), 1)

		inputFileData, err := ioutil.ReadFile(inputPath)
		if err != nil {
			panic(fmt.Errorf("error reading input file, %s", err))
		}
		userIdStrs := strings.Split(string(inputFileData), ",")
		userIds := make([]int64, 0, len(userIdStrs))
		for _, userIdStr := range userIdStrs {
			userId, err := strconv.ParseInt(userIdStr, 10, 64)
			if err != nil {
				panic(fmt.Errorf("error parsing userid %s to int64, %s", userIdStr, err))
			}
			userIds = append(userIds, userId)
		}

		fmt.Println("Starting GetUsersByIds")

		cfg := config.GetConfig()
		client := twitter.NewClient(cfg.APIKey, cfg.APISecret, cfg.Token, cfg.TokenSecret)
		users, err := client.GetUsersByIds(userIds)
		if err != nil {
			panic(err)
		}

		if len(users) > 0 {
			fmt.Printf("Start writing %d users to file %s\n", len(users), outputPath)
			file, err := os.OpenFile(outputPath, os.O_CREATE, 0666)
			defer file.Close()
			if err != nil {
				panic("error opening output file for write")
			}
			jsonEnc := json.NewEncoder(file)
			jsonEnc.SetIndent("", "    ")
			jsonEnc.Encode(users)
			fmt.Printf("Written %d users to file %s\n", len(users), outputPath)
		} else {
			fmt.Println("no users to write")
		}
	},
}

func init() {
	usersCmd.Flags().StringP("input", "i", "", "path to input user ids file")
	usersCmd.Flags().StringP("output", "o", "users_of_{input_file_name}_at_{date_time}.json", "path to output file")

	RootCmd.AddCommand(usersCmd)
}
