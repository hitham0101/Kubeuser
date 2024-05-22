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

type Contexts struct {
	Contexts       []string `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
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
			fmt.Println("Error creating .kubeuser directory", err)
			return
		}
		fmt.Println(".kubeuser directory created successfully")
	}

	// check if the cluster subdirectory already exists
	if _, err := os.Stat(filepath.Join(dirPath, clusterName)); !os.IsNotExist(err) {
		fmt.Printf("%s already initialized \n", clusterName)
		return
	} else {
		clusterSubDirectory, err := createClusterSubDirectory(clusterName, dirPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

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
		err = writeConfigToFile(config, clusterSubDirectory)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Configuration file created successfully.")
	}

	// ############################################################

	// Define the context file
	filePath := dirPath + "/context"

	// check if the context file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {

		addContextAndSetCurrent(filePath, clusterName)
		fmt.Println("context file updated successfully!")
		return
	} else {
		err = createContext(filePath, clusterName)
		if err != nil {
			log.Fatalf("Failed to create context file: %v", err)
		}
		fmt.Println("context file created successfully!")
	}

}

func createClusterSubDirectory(clusterName, kubeUserDir string) (dirPath string, err error) {

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

func createContext(filePath, clusterName string) error {

	// Define the configuration
	config := Contexts{
		Contexts:       []string{clusterName},
		CurrentContext: clusterName,
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

func addContextAndSetCurrent(filePath, newContext string) error {
	// Read the existing YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Parse the YAML file into the Contexts struct
	var Contexts Contexts
	err = yaml.Unmarshal(data, &Contexts)
	if err != nil {
		return fmt.Errorf("error unmarshalling YAML: %v", err)
	}

	// Add the new context to the contexts slice if it doesn't already exist
	for _, context := range Contexts.Contexts {
		if context == newContext {
			return fmt.Errorf("context '%s' already exists", newContext)
		}
	}
	Contexts.Contexts = append(Contexts.Contexts, newContext)

	// Set the new context as the current context
	Contexts.CurrentContext = newContext

	// Marshal the struct back into YAML
	updatedData, err := yaml.Marshal(&Contexts)
	if err != nil {
		return fmt.Errorf("error marshalling YAML: %v", err)
	}

	// Write the updated YAML back to the file
	err = os.WriteFile(filePath, updatedData, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	return nil
}

// func to check if there is a context file in the .kubeuser directory

func CheckContextFile() bool {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return false
	}

	homeDir := currentUser.HomeDir

	// Define the directory path
	dirPath := filepath.Join(homeDir, ".kubeuser")

	// Define the context file
	filePath := dirPath + "/context"

	// check if the context file already exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

// func to get the current context

func GetCurrentContext() string {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	homeDir := currentUser.HomeDir

	// Define the directory path
	dirPath := filepath.Join(homeDir, ".kubeuser")

	// Define the context file
	filePath := dirPath + "/context"

	// Read the existing YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	// Parse the YAML file into the Contexts struct
	var Contexts Contexts
	err = yaml.Unmarshal(data, &Contexts)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return ""
	}

	return Contexts.CurrentContext
}

// func to get the value of the current context in the subdirectory with the same name and had config.yaml file

func GetContextConfig(currentContext string) Config {
	// Get the current user's home directory
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return Config{}
	}

	homeDir := currentUser.HomeDir

	// Define the directory path
	dirPath := filepath.Join(homeDir, ".kubeuser")

	// Define the context file
	configFilePath := dirPath + "/" + currentContext + "/config.yaml"

	// Read the existing YAML file
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return Config{}
	}

	// Parse the YAML file into the Config struct
	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return Config{}
	}

	return config
}
