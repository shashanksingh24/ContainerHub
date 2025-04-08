package main

import (
	"fmt"
	"log"
	"os"

	"github.com/shashanksingh24/ContainerHub/pkg/rpc"

	"github.com/spf13/cobra"
)

var (
	socketPath = "/tmp/containerhub.sock"
)

func main() {
	rootCmd := &cobra.Command{Use: "containerhub"}

	createCmd := &cobra.Command{
		Use:   "create [image] [name] [command]",
		Short: "Create a new container",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			id, err := client.CreateContainer(args[0], args[1], args[2])
			if err != nil {
				log.Fatalf("Create failed: %v", err)
			}
			fmt.Printf("Container created with ID: %s\n", id)
		},
	}

	startCmd := &cobra.Command{
		Use:   "start [container_id]",
		Short: "Start a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			success, err := client.StartContainer(args[0])
			if err != nil {
				log.Fatalf("Start failed: %v", err)
			}
			if success {
				fmt.Println("Container started successfully")
			} else {
				fmt.Println("Failed to start container")
			}
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop [container_id]",
		Short: "Stop a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			success, err := client.StopContainer(args[0])
			if err != nil {
				log.Fatalf("Stop failed: %v", err)
			}
			if success {
				fmt.Println("Container stopped successfully")
			} else {
				fmt.Println("Failed to stop container")
			}
		},
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [container_id]",
		Short: "Delete a stopped container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			success, err := client.DeleteContainer(args[0])
			if err != nil {
				log.Fatalf("Delete failed: %v", err)
			}
			if success {
				fmt.Println("Container deleted successfully")
			} else {
				fmt.Println("Failed to delete container (ensure it’s stopped)")
			}
		},
	}

	execCmd := &cobra.Command{
		Use:   "exec [container_id] [command]",
		Short: "Execute a command in a running container",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			output, success, err := client.ExecCommand(args[0], args[1])
			if err != nil {
				log.Fatalf("Exec failed: %v", err)
			}
			if success {
				fmt.Printf("Output:\n%s\n", output)
			} else {
				fmt.Println("Failed to execute command")
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all containers",
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			containers, err := client.ListContainers()
			if err != nil {
				log.Fatalf("List failed: %v", err)
			}
			for _, c := range containers {
				fmt.Printf("ID: %s | Name: %s | Image: %s | Status: %s\n", c.Id, c.Name, c.Image, c.Status)
			}
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [container_id]",
		Short: "Retrieve logs from a container",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			client := rpc.NewClient(socketPath)
			defer client.Close()
			logs, success, err := client.GetContainerLogs(args[0])
			if err != nil {
				log.Fatalf("Logs retrieval failed: %v", err)
			}
			if success {
				fmt.Printf("Logs:\n%s\n", logs)
			} else {
				fmt.Println("Failed to retrieve logs")
			}
		},
	}

	rootCmd.AddCommand(createCmd, startCmd, stopCmd, deleteCmd, execCmd, listCmd, logsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
