package main

import (
	"context"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	greetpb "grpc-course/greet/greet_pb"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Infof("Processing unary request: %v", req)

	firstName := req.Greeting.FirstName
	lastName := req.Greeting.LastName
	result := "Hello " + firstName + " " + lastName

	return &greetpb.GreetResponse{
		Result: result,
	}, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Infof("Processing streaming request: %v", req)

	firstName := req.Greeting.FirstName
	lastName := req.Greeting.LastName

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " " + lastName + ". " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)

		time.Sleep(time.Second)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Info("LongGreet was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error reading from stream: %v", err)
			return err
		}

		firstName := req.Greeting.FirstName
		result += firstName
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Infof("Processing streaming request for bi dir streaming")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		firstName := req.Greeting.FirstName
		res := "Hello! " + firstName + "!"
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: res,
		})
		if err != nil {
			log.Fatalf("Error while sending request: %v", err)
			return err
		}
	}
}

func main() {
	log.Infof("Setting up server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("error listening: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error serving: %s", err)
	}
}
