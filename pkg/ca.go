package pkg

import (
	"context"
	"fmt"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func FetchCa(user, private_key_path, master_ip string) {

	master_socket := master_ip + ":22"

	// Use SSH key authentication from the auth package
	clientConfig, _ := auth.PrivateKey(user, private_key_path, ssh.InsecureIgnoreHostKey())

	// Create a new SCP client
	client := scp.NewClient(master_socket, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		fmt.Println("Couldn't establish a connection to the remote server ", err)
		return
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// ***************************************************************************************

	// Copy ca.crt from remote master machine to local machine

	// Create a local file to write to.
	f, err := os.OpenFile("./ca.crt", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Couldn't create the ca.crt file")
	}

	defer f.Close()

	// Copy the file from the remote server to the local file
	err = client.CopyFromRemote(context.Background(), f, "/etc/kubernetes/ssl/ca.crt")
	if err != nil {
		fmt.Printf("Copy failed from remote: %s", err.Error())
		os.Remove("./ca.crt")
		os.Exit(1)
	}

	fmt.Println("ca.crt copied successfully")

	// ***************************************************************************************

	// Copy ca.key from remote master machine to local machine

	// Create a local file to write to.
	f, err = os.OpenFile("./ca.key", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Couldn't create the ca.key file")
	}

	defer f.Close()

	// Copy the file from the remote server to the local file
	err = client.CopyFromRemote(context.Background(), f, "/etc/kubernetes/ssl/ca.key")
	if err != nil {
		fmt.Printf("Copy failed from remote: %s", err.Error())
		os.Remove("./ca.key")
		os.Exit(1)
	}

	fmt.Println("ca.key copied successfully")

}
