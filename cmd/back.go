package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// backCmd represents the back command
var backCmd = &cobra.Command{
	Use:   "back",
	Short: "Excute backend stress test",
	Long:  `Start backend stress test.`,
	Run: func(cmd *cobra.Command, args []string) {
		StartBackendStressTest()
	},
}

func init() {
	rootCmd.AddCommand(backCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// backCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// backCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartBackendStressTest() {
	se, err := NewScenarioExecuter("back", backendScenario)
	if err != nil {
		log.Fatal("Failed to NewScenarioExecuter", err)
		return
	}
	se.StartScenario()
}

func backendScenario(sc *ScenarioContext) {

	accountId := calcLoginId(sc.thNum, sc.loopNum)

	res := sc.client.SampleBLoginAccount(accountId, "password")
	sc.saveResult(res)

	res = sc.client.SampleBCouponList("0000000000001")
	sc.saveResult(res)
	res = sc.client.SampleBSendReserveList()
	sc.saveResult(res)
	res = sc.client.SampleBRecommendsList()
	sc.saveResult(res)

	res = sc.client.SampleBAccountDetail(1)
	sc.saveResult(res)
	res = sc.client.SampleBStampManagerList()
	sc.saveResult(res)
	res = sc.client.SampleBDeliveryCouponsList()
	sc.saveResult(res)
	res = sc.client.SampleBLogUseList()
	sc.saveResult(res)
	res = sc.client.SampleBDocumentsList()
	sc.saveResult(res)
	res = sc.client.SampleBDocumentDetail(1)
	sc.saveResult(res)
	res = sc.client.SampleBDlSummarysList("2020-12-01", "2020-12-03")
	sc.saveResult(res)
	res = sc.client.SampleBActSummaryList("2020-12-01", "2020-12-03")
	sc.saveResult(res)
	res = sc.client.SampleBPositionSummarysList("2020-12-01", "2020-12-03", "2020-12-02")
	sc.saveResult(res)
}

func calcLoginId(thNum, loopNum int) string {
	return fmt.Sprintf("%d", 100000+(thNum-1)*loopNum+loopNum)
}
