package pkg

import (
	"fmt"
	"os/exec"
)

func CheckKubectl() {
	cmd := exec.Command("kubectl", "version")
	if err := cmd.Run(); err != nil {
		fmt.Println("kubectl is not installed or not in the system PATH.")
	} else {
		fmt.Println("kubectl is installed and available.")
	}
}

func SetCluster(clusterName, apiEndpoint string) {
	// Define the command to execute
	command := exec.Command("kubectl", "config", "set-cluster", clusterName,
		"--certificate-authority=ca.crt",
		"--embed-certs=true",
		"--server="+apiEndpoint, //  https://10.1.0.101:6443
		"--kubeconfig=config")

	// Run the command and check for errors
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output of the command
	fmt.Println(string(output))
}

func SetCredentials(userName string) {
	// Define the command to execute
	command := exec.Command("kubectl", "config", "set-credentials", userName,
		"--client-certificate="+userName+".crt",
		"--client-key="+userName+".key",
		"--embed-certs=true",
		"--kubeconfig=config")

	// Run the command and check for errors
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output of the command
	fmt.Println(string(output))
}

func SetContext(userName, clusterName string) {
	// Define the command to execute
	command := exec.Command("kubectl", "config", "set-context", "default",
		"--cluster="+clusterName,
		"--user="+userName,
		"--kubeconfig=config")

	// Run the command and check for errors
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output of the command
	fmt.Println(string(output))
}

func UseContext() {
	// Define the command to execute
	command := exec.Command("kubectl", "config", "use-context", "default",
		"--kubeconfig=config")

	// Run the command and check for errors
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the output of the command
	fmt.Println(string(output))
}
