package subcommands

import (
	"github.com/spf13/cobra"
	"github.com/tun-io/tun-io/client"
)

var clientConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a tun-io server",
	Long:  "Connect to a tun-io server and forward traffic from the local machine to the server.",
	Run: func(cmd *cobra.Command, args []string) {
		client.Connect()
	},
}

var clientRootCmd = &cobra.Command{
	Use:   "client",
	Short: "tun-io client",
	Long:  "tun-io client is used to connect to a tun-io server and forward traffic.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	clientRootCmd.AddCommand(clientConnectCmd)

	clientConnectCmd.Flags().StringP("to", "t", "localhost:8000", "Target server address to connect to (default is localhost:8000)")
}

func GetClientRootCmd() *cobra.Command {
	return clientRootCmd
}
