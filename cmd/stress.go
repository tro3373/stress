package cmd

import (
	"log"
	"sync"

	"github.com/spf13/cobra"
)

// stressCmd represents the hello command
var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "Excute stress test",
	Long:  `Start stress test.`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(stressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func start() {
	log.Println(">> stress test start")

	sc, err := config.GetScenarioConfig("stress")
	if err != nil {
		log.Fatal("Failed to get scenario config")
	}

	wg := &sync.WaitGroup{}

	thNum := 0
	for {
		thNum++
		wg.Add(1)
		se := NewScenarioExecuter(thNum, sc.Count, wg)
		log.Println(">>> Starting ", se.String())
		go se.Start()
		if thNum >= sc.Thread {
			break
		}
	}
	wg.Wait()

	log.Println(">> stress test done")
}
