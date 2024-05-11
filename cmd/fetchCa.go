package cmd

import (
	"log"

	"github.com/hitham0101/kubeuser/pkg"
	"github.com/spf13/cobra"
)

// FetchCaCmd represents the add command
var FetchCaCmd = &cobra.Command{
	Use:   "fetch-ca",
	Short: "Fetcha CA key and certificate from the k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {

		// extract the flags
		master_server_user, _ := cmd.Flags().GetString("master_server_user")
		private_key_path, _ := cmd.Flags().GetString("private_key_path")
		master_ip, _ := cmd.Flags().GetString("master_ip")

		// Call the FetchCa function from the pkg package to fetch the CA key and certificate from the k8s cluster
		pkg.FetchCa(master_server_user, private_key_path, master_ip)

	},
}

func init() {
	rootCmd.AddCommand(FetchCaCmd)

	// Add flags to the subcommand
	FetchCaCmd.Flags().String("master_server_user", "", "The username to be used to authenticate with the k8s cluster master server")
	FetchCaCmd.Flags().String("private_key_path", "", "The path to the private key file which will be used to authenticate with the k8s cluster master server")
	FetchCaCmd.Flags().String("master_ip", "", "The IP address of the k8s cluster master server")

	// Mark the flags as required flags for the subcommand
	err := FetchCaCmd.MarkFlagRequired("master_server_user")
	if err != nil {
		log.Fatal(err)
	}
	err = FetchCaCmd.MarkFlagRequired("private_key_path")
	if err != nil {
		log.Fatal(err)
	}
	err = FetchCaCmd.MarkFlagRequired("master_ip")
	if err != nil {
		log.Fatal(err)
	}

}
