package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// reguserCmd represents the reguser command
var reguserCmd = &cobra.Command{
	Use:   "reguser",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartRegistUsers()
	},
}

func init() {
	rootCmd.AddCommand(reguserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// reguserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// reguserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartRegistUsers() {
	se, err := NewScenarioExecuter("reguser", registUsers)
	if err != nil {
		log.Fatal("Failed to NewScenarioExecuter", err)
		return
	}
	se.StartScenario()
}

func registUsers(sc *ScenarioContext) {

	deviceType := "ios"
	if remainder := sc.loopNum % 2; remainder == 0 {
		deviceType = "adr"
	}
	res := sc.client.SampleFRegistDeviceId(1, deviceType, genDeviceId(sc))
	sc.saveResult(res)
}

func genDeviceId(sc *ScenarioContext) string {
	return fmt.Sprintf("sample_%s_%d_%010d", sc.startTime, sc.thNum, sc.loopNum)
}
