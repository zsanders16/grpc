package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials"

	"github.com/zsanders16/grpc/pb"
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
	pb.RegisterEmployeeServiceServer(s, new(employeeService))
	log.Println("starting server on port ", port)
	s.Serve(lis)
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(cxt context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	return nil, nil
}
func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {

	return nil
}
func (s *employeeService) Save(cxt context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {

	return nil, nil
}
func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {

	return nil
}
func (s *employeeService) AddPhoto(stram pb.EmployeeService_AddPhotoServer) error {

	return nil
}
