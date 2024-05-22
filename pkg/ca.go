package pkg

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func FetchCa(master_server_user, private_key_path, master_ip, ClusterName string) {

	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	homeDir := currentUser.HomeDir

	// Define the directory path
	clusterDirPath := filepath.Join(homeDir, ".kubeuser", ClusterName)

	master_socket := master_ip + ":22"

	// Use SSH key authentication from the auth package
	clientConfig, _ := auth.PrivateKey(master_server_user, private_key_path, ssh.InsecureIgnoreHostKey())

	// Create a new SCP client
	client := scp.NewClient(master_socket, &clientConfig)

	// Connect to the remote server
	err = client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// ***************************************************************************************

	// Copy ca.crt from remote master machine to local machine

	// create the ca directory
	err = os.MkdirAll(clusterDirPath+"/ca", 0755)
	if err != nil {
		fmt.Println("Couldn't create the ca directory")
	}
	// Create a local file to write to.

	f, err := os.OpenFile(clusterDirPath+"/ca/ca.crt", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Couldn't create the ca.crt file")
	}

	defer f.Close()

	// Copy the file from the remote server to the local file
	err = client.CopyFromRemote(context.Background(), f, "/etc/kubernetes/ssl/ca.crt")
	if err != nil {
		fmt.Printf("Copy failed from remote: %s", err.Error())
		os.Remove(clusterDirPath + "/ca/ca.crt")
		os.Exit(1)
	}

	fmt.Println("ca.crt has been downloaded successfully")

	// ***************************************************************************************

	// Copy ca.key from remote master machine to local machine

	// Create a local file to write to.
	f, err = os.OpenFile(clusterDirPath+"/ca/ca.key", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Couldn't create the ca.key file")
	}

	defer f.Close()

	// Copy the file from the remote server to the local file
	err = client.CopyFromRemote(context.Background(), f, "/etc/kubernetes/ssl/ca.key")
	if err != nil {
		fmt.Printf("Copy failed from remote: %s", err.Error())
		os.Remove(clusterDirPath + "/ca/ca.key")
		os.Exit(1)
	}

	fmt.Println("ca.key has been downloaded successfully")

}
