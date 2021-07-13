package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use: "time",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now().UnixNano() / int64(time.Millisecond)
		fmt.Printf("%d\n", t)
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
}
