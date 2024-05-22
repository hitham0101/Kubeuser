package cmd

import (
	"os"

	"github.com/hitham0101/kubeuser/pkg"
	"github.com/spf13/cobra"
)

// FetchCaCmd represents the add command
var FetchCaCmd = &cobra.Command{
	Use:   "fetch-ca",
	Short: "Fetcha CA key and certificate from the k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {

		// extract the flags
		private_key_path, _ := cmd.Flags().GetString("private_key_path")
		master_ip, _ := cmd.Flags().GetString("master_ip")
		master_server_user, _ := cmd.Flags().GetString("master_server_user")
		cluster_name, _ := cmd.Flags().GetString("cluster_name")

		if cluster_name == "" || master_ip == "" || master_server_user == "" || private_key_path == "" {

			var currentContext string = pkg.GetCurrentContext()
			var namConfige pkg.Config = pkg.GetContextConfig(currentContext)

			cluster_name = namConfige.ClusterName
			master_ip = namConfige.MasterIP
			master_server_user = namConfige.MasterServerUser
			private_key_path = namConfige.PrivateKeyPath
			pkg.FetchCa(master_server_user, private_key_path, master_ip, cluster_name)
			os.Exit(0)

		}

		// Call the FetchCa function from the pkg package to fetch the CA key and certificate from the k8s cluster
		pkg.FetchCa(master_server_user, private_key_path, master_ip, cluster_name)

	},
}

func init() {
	rootCmd.AddCommand(FetchCaCmd)

	// Add flags to the subcommand
	FetchCaCmd.Flags().String("master_server_user", "", "The username to be used to authenticate with the k8s cluster master server")
	FetchCaCmd.Flags().String("private_key_path", "", "The path to the private key file which will be used to authenticate with the k8s cluster master server")
	FetchCaCmd.Flags().String("master_ip", "", "The IP address of the k8s cluster master server")
	FetchCaCmd.Flags().String("cluster_name", "", "The name of the cluster")

}
