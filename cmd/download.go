/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a package version",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		downloadPackage(args[0], cmd.Flag("version").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	downloadCmd.Flags().StringP("version", "v", "latest", "package version")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func downloadPackage(packageName, versionName string) {
	url := fmt.Sprintf("%v/api/v1/packages/%v/versions/%v/", viper.GetString("server"), packageName, versionName)
	resp, err := Get(url)
	check(err)

	if resp.StatusCode != 200 {
		fmt.Println(resp.Error())
		return
	}
	dis := resp.Header.Get("Content-Disposition")
	fileName := fmt.Sprintf("%v_%v", packageName, versionName)
	if dis != "" {
		_, params, err := mime.ParseMediaType(dis)
		check(err)
		fileName = params["filename"]
	}
	f, err := os.Create(fileName)
	check(err)
	defer f.Close()

	io.Copy(f, resp.Body)
	defer resp.Body.Close()
	f.Sync()
}
