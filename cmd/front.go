package cmd

import (
	"log"

	"github.com/tro3373/stress/cmd/backend"

	"github.com/spf13/cobra"
)

// frontCmd represents the hello command
var frontCmd = &cobra.Command{
	Use:   "front",
	Short: "Excute frontend stress test",
	Long:  `Start frontend stress test.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartFrontendStressTest()
	},
}

func init() {
	rootCmd.AddCommand(frontCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// frontCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// frontCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartFrontendStressTest() {
	se, err := NewScenarioExecuter("front", frontendScenario)
	if err != nil {
		log.Fatal("Failed to NewScenarioExecuter", err)
		return
	}
	se.StartScenario()
}

func frontendScenario(sc *ScenarioContext) {
	var res *backend.Res
	res = sc.client.SampleFContentsList()
	sc.saveResult(res)
	res = sc.client.SampleFContentsDetail(1)
	sc.saveResult(res)
	res = sc.client.SampleFContentsDetail(2)
	sc.saveResult(res)
	res = sc.client.SampleFCouponList()
	sc.saveResult(res)
	res = sc.client.SampleFRecommendsList()
	sc.saveResult(res)
	res = sc.client.SampleFStampManagerList()
	sc.saveResult(res)

	res = sc.client.SampleFBrandsList()
	sc.saveResult(res)
	res = sc.client.SampleFShopList()
	sc.saveResult(res)
	res = sc.client.SampleFUserFavoriteBrandList()
	sc.saveResult(res)

	userId := "sample01"
	res = sc.client.SampleFDeliveryCouponsList(userId)
	sc.saveResult(res)
}
