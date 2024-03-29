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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	_package "github.com/syfun/package/pkg/package"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get package info",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		getPackage(args[0])
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// infoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getPackage(name string) {
	url := fmt.Sprintf("%v/api/v1/packages/%v/", viper.GetString("server"), name)
	resp, err := Get(url)
	check(err)
	if resp.StatusCode != 200 {
		fmt.Println(resp.Error())
		return
	}

	var p _package.Package
	check(resp.Decode(&p))

	fmt.Printf("Package:\n\nID: %v\nName: %v\n\n", p.ID, p.Name)
	if p.Versions == nil || len(p.Versions) == 0 {
		return
	}
	var data [][]interface{}
	for _, v := range p.Versions {
		data = append(data, []interface{}{v.ID, v.Name, v.Size, v.FileName})
	}
	printTable([]string{"ID", "Name", "Size", "FileName"}, data)
}