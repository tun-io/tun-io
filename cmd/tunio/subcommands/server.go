package subcommands

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/tun-io/tun-io/server"
)

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the tun-io server",
	Long:  "Start the tun-io server to accept incoming connections and forward traffic.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting the tun-io server...")
		server.StartServer()
	},
}

var serverRootCmd = &cobra.Command{
	Use:   "server",
	Short: "tun-io server",
	Long:  "tun-io is the primary server.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	serverRootCmd.AddCommand(serverStartCmd)
}

func GetServerRootCmd() *cobra.Command {
	return serverRootCmd
}
