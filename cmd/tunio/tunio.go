package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tun-io/tun-io/cmd/tunio/subcommands"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "tun-io",
	Short: "tun-io is a tunneling tool",
	Long:  "tun-io is a tunneling tool for testing websites on other devices without any port forwarding.",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("tun-io is a tunneling tool for testing websites on other devices without any port forwarding.\n")
	},
}

func init() {
	rootCmd.AddCommand(subcommands.GetClientRootCmd())
	rootCmd.AddCommand(subcommands.GetServerRootCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
