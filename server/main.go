package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc/metadata"

	"github.com/zsanders16/grpc/pb"
	"github.com/zsanders16/grpc/server/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const port = ":9000"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile("../cert.pem", "../key.pem")
	if err != nil {
		log.Fatal(err)
	}
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	s := grpc.NewServer(opts...)
	pb.RegisterEmployeeServiceServer(s, &employeeService{})
	log.Println("starting server on port ", port)
	s.Serve(lis)
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("Metadata received: %v\n", md)
	}

	for _, e := range data.Employees {
		if req.BadgeNumber == e.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {

	for _, e := range data.Employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}

	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Badge Number: %v\n", md["badgenumber"][0])
	}

	imgData := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Recieved file with length: %v\n", len(imgData))
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Recieved %v bytes\n", len(data.Data))
		imgData = append(imgData, data.Data...)
	}
}

func (s *employeeService) Save(cxt context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {

	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {

	return nil
}
