/*
Copyright 息 2020 NAME HERE <EMAIL ADDRESS>

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
	// "bytes"
	// "encoding/json"
	"fmt"
	"os"

	// "io/ioutil"
	// "net/http"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

// helloCmd represents the hello command
var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helloCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helloCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func start() {
	fmt.Println(">> hello called")
	// fmt.Printf(">>> %#v\n", config)
	res, err := Req(ContentsDetail)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("> res %s.\n", res)

	// // fmt.Println(string(byteArray))
	// // fmt.Println(Pretty(string(byteArray), ""))
	// var buf bytes.Buffer
	// err = json.Indent(&buf, []byte(res.json), "", "  ")
	// if err != nil {
	// 	fmt.Println(">> Failed to parse json", err)
	// 	os.Exit(1)
	// 	// return HandleReqError("parse json", err)
	// }
	// fmt.Println(">> Res:", buf.String())

	// url := "https://dev.app.tenco.co.jp/contents-api/v1/frontend/contents/39"
	// req, _ := http.NewRequest("GET", url, nil)
	// req.Header.Set("sample-header", "sample-value")
	// client := new(http.Client)
	// resp, _ := client.Do(req)
	// defer resp.Body.Close()
	// byteArray, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(byteArray))
}
