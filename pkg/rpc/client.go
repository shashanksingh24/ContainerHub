package rpc

import (
	"context"
	"log"

	pb "github.com/shashanksingh24/ContainerHub/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn *grpc.ClientConn
	svc  pb.ContainerServiceClient
}

func NewClient(socketPath string) *Client {
	conn, err := grpc.Dial("unix://"+socketPath, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	return &Client{
		conn: conn,
		svc:  pb.NewContainerServiceClient(conn),
	}
}

func (c *Client) CreateContainer(image, name, command string) (string, error) {
	resp, err := c.svc.CreateContainer(context.Background(), &pb.CreateRequest{
		Image:   image,
		Name:    name,
		Command: command,
	})
	if err != nil {
		return "", err
	}
	return resp.ContainerId, nil
}

func (c *Client) StartContainer(containerID string) (bool, error) {
	resp, err := c.svc.StartContainer(context.Background(), &pb.StartRequest{ContainerId: containerID})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (c *Client) StopContainer(containerID string) (bool, error) {
	resp, err := c.svc.StopContainer(context.Background(), &pb.StopRequest{ContainerId: containerID})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (c *Client) DeleteContainer(containerID string) (bool, error) {
	resp, err := c.svc.DeleteContainer(context.Background(), &pb.DeleteRequest{ContainerId: containerID})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

func (c *Client) ExecCommand(containerID, command string) (string, bool, error) {
	resp, err := c.svc.ExecCommand(context.Background(), &pb.ExecRequest{ContainerId: containerID, Command: command})
	if err != nil {
		return "", false, err
	}
	return resp.Output, resp.Success, nil
}

func (c *Client) ListContainers() ([]*pb.ContainerInfo, error) {
	resp, err := c.svc.ListContainers(context.Background(), &pb.ListRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Containers, nil
}

func (c *Client) GetContainerLogs(containerID string) (string, bool, error) {
	resp, err := c.svc.GetContainerLogs(context.Background(), &pb.LogsRequest{ContainerId: containerID})
	if err != nil {
		return "", false, err
	}
	return resp.Logs, resp.Success, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
