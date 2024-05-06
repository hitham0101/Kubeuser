# kubeuser CLI Tool

## Overview
kubeuser CLI Tool is a command-line interface utility designed to manage K8S users efficiently. With this tool, you can easily add, remove and list users on your kubernetes cluster

## Features
- Add new users
- Remove existing users #TODO
- List all users #TODO

## Build
1. Clone this repository to your local machine.
2. Navigate to the directory containing the cloned repository.
3. Run the following command to build the CLI tool:
    ```
    make build
    ```
    This will create a binary file named `kubeuser` in the `bin` directory.
    Note: You should have Go installed on your machine to build the CLI tool.


## Installation
1. Move the `kubeuser` binary to a directory in your system path.
    ```
    mv bin/kubeuser /usr/local/bin
    ```

2. Make the binary executable.
    ```
    chmod +x /usr/local/bin/kubeuser
    ```

## kubeuser Prerequisites

1. You should have a kubectl  installed on your machine. 
2. You should have ca.key and ca.crt for your kubernetes cluster


## Usage


### Adding a New User
To add a new user, use the following command:

```
kubeuser add --user_name <username>
```

### Removing an Existing User
To remove an existing user, use the following command:

```
kubeuser remove --user_name <username>
```

### Listing All  Users

To list all users, use the following command:

```
kubeuser list
```


