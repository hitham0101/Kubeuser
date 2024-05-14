package pkg

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClusterName      string `yaml:"cluster_name"`
	PrivateKeyPath   string `yaml:"private_key_path"`
	MasterServerUser string `yaml:"master_server_user"`
	MasterIP         string `yaml:"master_ip"`
}

// Config represents the structure of the YAML file
type Context struct {
	Context string `yaml:"context"`
}

func Initialize(cluster_name, master_ip, master_server_user, private_key_path string) {

	clusterName := cluster_name
	privateKeyPath := private_key_path
	masterServerUser := master_server_user
	masterIP := master_ip

	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	homeDir := currentUser.HomeDir

	// Define the directory path
	dirPath := filepath.Join(homeDir, ".kubeuser")

	// Check if the directory already exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// Create the directory
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			fmt.Println("Error creating .kubeuser directory:", err)
			return
		}
		fmt.Println(".kubeuser directory created successfully:", dirPath)
	} else {
		fmt.Println(".kubeuser directory already exists:", dirPath)
	}

	subdirectory, err := createSubDirectory(clusterName, dirPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Subdirectory created successfully.")

	// ############################################################
	// Write configuration to file

	// Create configuration object
	config := Config{
		ClusterName:      clusterName,
		PrivateKeyPath:   privateKeyPath,
		MasterServerUser: masterServerUser,
		MasterIP:         masterIP,
	}

	// Write configuration to file
	err = writeConfigToFile(config, subdirectory)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Configuration file created successfully.")

	filePath := dirPath + "/context"
	err = createContext(filePath)
	if err != nil {
		log.Fatalf("Failed to create context file: %v", err)
	}

	fmt.Println("context file created successfully!")

}

func createSubDirectory(clusterName, kubeUserDir string) (dirPath string, err error) {

	// Define the subdirectory path
	clusterDirPath := filepath.Join(kubeUserDir, clusterName)

	// Check if the .kubeuser directory exists
	if _, err := os.Stat(kubeUserDir); os.IsNotExist(err) {
		// If .kubeuser directory doesn't exist, create it
		if err := os.Mkdir(kubeUserDir, 0755); err != nil {
			return "", err
		}
	}

	// Check if the subdirectory already exists
	if _, err := os.Stat(clusterDirPath); os.IsNotExist(err) {
		// Create the subdirectory
		if err := os.Mkdir(clusterDirPath, 0755); err != nil {
			return "", err
		}
	} else {
		fmt.Println("Subdirectory already exists:", clusterDirPath)
	}

	return clusterDirPath, nil
}

func writeConfigToFile(config Config, dirPath string) error {

	path := filepath.Join(dirPath, "config.yaml")
	// Marshal configuration object to YAML
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	// Open file for writing
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write YAML data to file
	_, err = file.Write(yamlData)
	if err != nil {
		return err
	}

	return nil
}

func createContext(filePath string) error {

	fmt.Println("Creating context file...")
	fmt.Println("Context file path:", filePath)
	// Define the configuration
	config := Context{
		Context: "default",
	}

	// Marshal the config struct to YAML format
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("error marshalling config to YAML: %w", err)
	}

	// Write the YAML data to the specified file
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %w", err)
	}

	return nil
}
