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
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload a package version",
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		uploadPackage(args[0], cmd.Flag("version").Value.String(), args[1])
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	uploadCmd.Flags().StringP("version", "v", "latest", "package name")
}

func uploadPackage(packageName, versionName, filePath string) {
	buf := new(bytes.Buffer)
	w := multipart.NewWriter(buf)
	if versionName == "" {
		versionName = "latest"
	}

	if err := w.WriteField("name", versionName); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fw, err := w.CreateFormFile("file", fileInfo.Name())
	if _, err := io.Copy(fw, file); err != nil {
		log.Fatal(err)
	}
	w.Close()

	url := fmt.Sprintf("%v/api/v1/packages/%v/versions/", viper.GetString("server"), packageName)
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		r := Response{resp}
		fmt.Println(r.Error())
		return
	}
	fmt.Println("Added package version")
}
