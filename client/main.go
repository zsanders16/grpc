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
	case 4:
		SaveEmployee(client)
	case 5:
		SaveAll(client)
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

func SaveEmployee(client pb.EmployeeServiceClient) {
	e := &pb.Employee{
		Id:                  4,
		BadgeNumber:         6358,
		FirstName:           "John",
		LastName:            "Smith",
		VacationAccrualRate: 3,
		VacationAccred:      33,
	}
	d, err := client.Save(context.Background(), &pb.EmployeeRequest{Employee: e})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(d)
}

func SaveAll(client pb.EmployeeServiceClient) {
	newEmps := []pb.Employee{
		pb.Employee{
			BadgeNumber:         8264,
			FirstName:           "Chris",
			LastName:            "Adams",
			VacationAccrualRate: 1.3,
			VacationAccred:      7,
		},
		pb.Employee{
			BadgeNumber:         9361,
			FirstName:           "Amy",
			LastName:            "Johns",
			VacationAccrualRate: 2.4,
			VacationAccred:      14,
		},
		pb.Employee{
			BadgeNumber:         1743,
			FirstName:           "George",
			LastName:            "Mack",
			VacationAccrualRate: 3,
			VacationAccred:      42,
		},
	}

	stream, err := client.SaveAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	doneCh := make(chan bool)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				doneCh <- true
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(res.Employee)
		}
	}()

	for _, e := range newEmps {
		err := stream.Send(&pb.EmployeeRequest{Employee: &e})
		if err != nil {
			log.Fatal(err)
		}
	}
	stream.CloseSend()
	<-doneCh
}
