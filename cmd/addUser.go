package cmd

import (
	"log"

	"github.com/hitham0101/kubeuser/pkg"
	"github.com/spf13/cobra"
)

// AddUserCmd represents the add command
var AddUserCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a user to the k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {

		// extract the flags
		user, _ := cmd.Flags().GetString("user_name")
		cluster, _ := cmd.Flags().GetString("cluster_name")
		api_server_url, _ := cmd.Flags().GetString("api_server_url")

		pkg.GeneratePrivateKey(user)
		pkg.GenerateCSR(user)
		pkg.GenerateCertificate(user)
		pkg.CheckKubectl()
		pkg.SetCluster(cluster, api_server_url) 
		pkg.SetCredentials(user)
		pkg.SetContext(user, cluster)
		pkg.UseContext()

	},
}

func init() {
	rootCmd.AddCommand(AddUserCmd)

	// Add flags to the subcommand
	AddUserCmd.Flags().String("user_name", "", "The name of the user to be added to the k8s cluster")
	AddUserCmd.Flags().String("cluster_name", "", "The name of the cluster")
	AddUserCmd.Flags().String("api_server_url", "", "The name of the cluster")

	// Mark the flags as required flags for the subcommand
	err := AddUserCmd.MarkFlagRequired("user_name")
	if err != nil {
		log.Fatal(err)
	}
	err = AddUserCmd.MarkFlagRequired("cluster_name")
	if err != nil {
		log.Fatal(err)
	}

	err = AddUserCmd.MarkFlagRequired("api_server_url")
	if err != nil {
		log.Fatal(err)
	}

}
