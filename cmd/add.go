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
	"log"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add package",
	Long: "",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addPackage(args[0])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	//addCmd.SetUsageTemplate("pkg add NAME")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//addCmd.Flags().StringP("name", "n", "", "package name(required)")
	//addCmd.MarkFlagRequired("name")
}

func addPackage(name string) {
	_, err := Post(viper.GetString("server") + "/api/v1/packages/", JSON{"name": name})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Added %v\n", name)
}