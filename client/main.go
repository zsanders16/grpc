package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/zsanders16/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

const port = ":9000"

func main() {
	option := flag.Int("o", 1, "Command to run")
	flag.Parse()
	creds, err := credentials.NewClientTLSFromFile("../cert.pem", "")
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	conn, err := grpc.Dial("localhost"+port, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := pb.NewEmployeeServiceClient(conn)
	switch *option {
	case 1:
		GetEmp2080(client)
	case 2:
		GetAllEmp(client)
	case 3:
		SendPhoto(client)
	}

}

func GetEmp2080(client pb.EmployeeServiceClient) {

	md := metadata.MD{}
	md["user"] = []string{"zsanders"}
	md["password"] = []string{"password1"}

	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	e, err := client.GetByBadgeNumber(ctx, &pb.GetByBadgeNumberRequest{BadgeNumber: 2080})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(e)
}

func GetAllEmp(client pb.EmployeeServiceClient) {

	stream, err := client.GetAll(context.Background(), &pb.GetAllRequest{})
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res.Employee)
	}
}

func SendPhoto(client pb.EmployeeServiceClient) {
	f, err := os.Open("../Penguins.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	md := metadata.New(map[string]string{
		"badgenumber": "2080",
	})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err := client.AddPhoto(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for {
		chunk := make([]byte, 64*1024)
		n, err := f.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if n < len(chunk) {
			chunk = chunk[:n]
		}
		stream.Send(&pb.AddPhotoRequest{Data: chunk})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.IsOk)
}
