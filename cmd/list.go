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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("list called")
		listPackages(cmd.Flag("fuzzyName").Value.String())
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().StringP("fuzzyName", "n", "", "fuzzy name")
}

func printTable(headers []string, data [][]interface{}) {
	l := len(headers)
	if l == 0 {
		return
	}

	s := ""
	for _, h := range headers {
		s += h + "\t"
	}
	fmt.Println(s)

	for _, d := range data {
		s := ""
		for i, v := range d {
			s = fmt.Sprintf("%v%v\t", s, v)
			if i+1 > l {
				break
			}
		}
		fmt.Println(s)
	}
}

func listPackages(fuzzyName string) {
	url := fmt.Sprintf(viper.GetString("server") + "/api/v1/packages/?fuzzy_nam=%v", fuzzyName)
	resp, err := Get(url)
	check(err)
	if resp.StatusCode != 200 {
		fmt.Println(resp.Error())
		return
	}

	var packages []_package.Package
	check(resp.Decode(&packages))

	var data [][]interface{}
	for _, p := range packages {
		data = append(data, []interface{}{p.ID, p.Name})
	}
	printTable([]string{"ID", "Name"}, data)
}
