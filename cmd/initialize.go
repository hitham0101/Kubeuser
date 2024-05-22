package cmd

import (
	"log"

	"github.com/hitham0101/kubeuser/pkg"

	"github.com/spf13/cobra"
)

// initCmd represents the add command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init kubeuser tool with basic configurations",
	Run: func(cmd *cobra.Command, args []string) {

		// extract the flags
		cluster_name, _ := cmd.Flags().GetString("cluster_name")
		master_ip, _ := cmd.Flags().GetString("master_ip")
		master_server_user, _ := cmd.Flags().GetString("master_server_user")
		private_key_path, _ := cmd.Flags().GetString("private_key_path")

		pkg.Initialize(cluster_name, master_ip, master_server_user, private_key_path)

	},
}

func init() {

	rootCmd.AddCommand(initCmd)

	// Add flags to the subcommand
	initCmd.Flags().String("cluster_name", "", "The name of the cluster")
	initCmd.Flags().String("master_ip", "", "The IP address of the k8s cluster master server")
	initCmd.Flags().String("master_server_user", "", "The username to be used to authenticate with the k8s cluster master server")
	initCmd.Flags().String("private_key_path", "", "The path to the private key file which will be used to authenticate with the k8s cluster master server")

	// Mark the flags as required flags for the subcommand

	err := initCmd.MarkFlagRequired("cluster_name")
	if err != nil {
		log.Fatal(err)
	}

	err = initCmd.MarkFlagRequired("master_ip")
	if err != nil {
		log.Fatal(err)
	}

}
