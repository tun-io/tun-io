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
		remote, _ := cmd.Flags().GetString("remote")
		local, _ := cmd.Flags().GetString("local")
		secure, _ := cmd.Flags().GetBool("secure")

		client.SetupStore(remote, local, secure)

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

	clientConnectCmd.Flags().StringP("remote", "r", "a.tunio.test", "The remote domain to connect to")
	clientConnectCmd.Flags().StringP("local", "l", "localhost:8000", "The local domain to forward traffic to")
	clientConnectCmd.Flags().BoolP("secure", "s", false, "Use secure connection (wss)")
}

func GetClientRootCmd() *cobra.Command {
	return clientRootCmd
}
