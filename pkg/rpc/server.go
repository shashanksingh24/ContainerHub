package rpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	pb "github.com/shashanksingh24/ContainerHub/proto"
	"github.com/shashanksingh24/ContainerHub/pkg/container"

	"google.golang.org/grpc"
	"math/rand"
)

type Server struct {
	pb.UnimplementedContainerServiceServer
	mu         sync.Mutex
	containers map[string]*container.Container
}

func NewServer() *Server {
	return &Server{
		containers: make(map[string]*container.Container),
	}
}

func (s *Server) CreateContainer(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := &container.Container{
		ID:      generateID(),
		Name:    req.Name,
		Image:   req.Image,
		Command: req.Command,
		Status:  "created",
	}
	s.containers[c.ID] = c

	err := c.PrepareBundle()
	if err != nil {
		log.Printf("Failed to prepare bundle for %s: %v", c.ID, err)
		return nil, err
	}
	log.Printf("Container %s created", c.ID)
	return &pb.CreateResponse{ContainerId: c.ID}, nil
}

func (s *Server) StartContainer(ctx context.Context, req *pb.StartRequest) (*pb.StartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, exists := s.containers[req.ContainerId]
	if !exists {
		log.Printf("Container %s not found", req.ContainerId)
		return &pb.StartResponse{Success: false}, nil
	}

	err := exec.Command("runc", "run", "-d", "--log", filepath.Join("/tmp/containerhub", c.ID, "log.json"), c.ID).Run()
	if err != nil {
		log.Printf("Failed to start %s: %v", c.ID, err)
		return nil, err
	}
	c.Status = "running"
	log.Printf("Container %s started", c.ID)
	return &pb.StartResponse{Success: true}, nil
}

func (s *Server) StopContainer(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, exists := s.containers[req.ContainerId]
	if !exists {
		log.Printf("Container %s not found", req.ContainerId)
		return &pb.StopResponse{Success: false}, nil
	}

	err := exec.Command("runc", "kill", c.ID, "SIGTERM").Run()
	if err != nil {
		log.Printf("Failed to stop %s: %v", c.ID, err)
		return nil, err
	}
	c.Status = "stopped"
	log.Printf("Container %s stopped", c.ID)
	return &pb.StopResponse{Success: true}, nil
}

func (s *Server) DeleteContainer(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, exists := s.containers[req.ContainerId]
	if !exists || c.Status == "running" {
		log.Printf("Cannot delete %s: not found or still running", req.ContainerId)
		return &pb.DeleteResponse{Success: false}, nil
	}

	err := os.RemoveAll(filepath.Join("/tmp/containerhub", c.ID))
	if err != nil {
		log.Printf("Failed to delete bundle for %s: %v", c.ID, err)
		return nil, err
	}
	delete(s.containers, c.ID)
	log.Printf("Container %s deleted", c.ID)
	return &pb.DeleteResponse{Success: true}, nil
}

func (s *Server) ExecCommand(ctx context.Context, req *pb.ExecRequest) (*pb.ExecResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, exists := s.containers[req.ContainerId]
	if !exists || c.Status != "running" {
		log.Printf("Cannot exec in %s: not found or not running", req.ContainerId)
		return &pb.ExecResponse{Success: false}, nil
	}

	out, err := exec.Command("runc", "exec", c.ID, "sh", "-c", req.Command).CombinedOutput()
	if err != nil {
		log.Printf("Exec failed in %s: %v", c.ID, err)
		return &pb.ExecResponse{Output: string(out), Success: false}, nil
	}
	log.Printf("Executed command in %s: %s", c.ID, req.Command)
	return &pb.ExecResponse{Output: string(out), Success: true}, nil
}

func (s *Server) ListContainers(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var containers []*pb.ContainerInfo
	for _, c := range s.containers {
		containers = append(containers, &pb.ContainerInfo{
			Id:     c.ID,
			Name:   c.Name,
			Image:  c.Image,
			Status: c.Status,
		})
	}
	log.Printf("Listed %d containers", len(containers))
	return &pb.ListResponse{Containers: containers}, nil
}

func (s *Server) GetContainerLogs(ctx context.Context, req *pb.LogsRequest) (*pb.LogsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, exists := s.containers[req.ContainerId]
	if !exists {
		log.Printf("Container %s not found", req.ContainerId)
		return &pb.LogsResponse{Success: false}, nil
	}

	logFile := filepath.Join("/tmp/containerhub", c.CID, "log.json")
	logs, err := os.ReadFile(logFile)
	if err != nil {
		log.Printf("Failed to read logs for %s: %v", c.ID, err)
		return &pb.LogsResponse{Logs: "No logs available", Success: false}, nil
	}
	log.Printf("Retrieved logs for %s", c.ID)
	return &pb.LogsResponse{Logs: string(logs), Success: true}, nil
}

func generateID() string {
	return "hub_" + fmt.Sprintf("%08d", rand.Intn(100000000))
}

func StartServer(socketPath string) {
	lis, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterContainerServiceServer(s, NewServer())
	log.Printf("ContainerHub server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
