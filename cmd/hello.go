package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	// "reflect"

	"github.com/spf13/cobra"
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
	log.Println(">> hello called")
	client := NewApiClient(config)
	res, err := client.GetContentsDetail()
	if err != nil {
		log.Fatal(err)
	}
	write(1, res.ReqNo, res.Json)
	// write(3, res.ReqNo, res.Json)
}

func write(thNum, rqNum int, data string) {
	path := outPath(config, 1, rqNum)
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, data)
}

func outPath(config Config, thNum, rqNum int) string {
	path := filepath.Join(config.LogDir, fmt.Sprintf("ThreadNo_%03d", thNum), fmt.Sprintf("ReqNo_%03d.txt", rqNum))
	return path
}
