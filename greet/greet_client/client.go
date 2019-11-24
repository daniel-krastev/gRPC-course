package main

import (
	"context"
	"io"
	"os"
	"time"

	greetpb "grpc-course/greet/greet_pb"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("Setting up client...")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error dialing: %s", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	// doUnaryCall(c)

	// doServerStreamCall(c)

	// doClientStreamCall(c)

	doBiDiStreaming(c)
}

func doUnaryCall(c greetpb.GreetServiceClient) {
	log.Info("Calling greet server unary call...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daniel",
			LastName:  "Krastev",
		},
	}
	resp, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet unary service: %v", err)
	}

	log.Infof("Response from greet unary call: %s", resp)
}

func doServerStreamCall(c greetpb.GreetServiceClient) {
	log.Info("Calling greet server server streaming call...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Daniel",
			LastName:  "Krastev",
		},
	}
	resultStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling greet server streaming service: %v", err)
	}
	for {
		resp, err := resultStream.Recv()
		if err == io.EOF {
			log.Info("Reached end of stream for result of call to server streaming...")
			break
		}
		if err != nil {
			log.Fatalf("Error while calling greet server stream service: %v", err)
		}
		log.Infof("Response from greet server stream call: %s", resp)
	}
}

func doClientStreamCall(c greetpb.GreetServiceClient) {
	log.Info("Calling greet server client streaming call...")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling greet server client streaming service: %v", err)
	}
	reqs := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani",
				LastName:  "Kras",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani2",
				LastName:  "Kras2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani3",
				LastName:  "Kras3",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani4",
				LastName:  "Kras4",
			},
		},
	}

	for _, req := range reqs {
		log.Infof("Sending req to long greet with %v", req)
		stream.SendMsg(req)
		time.Sleep(time.Second * 2)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Received err from server while calling long greet request: %v", err)
	}
	log.Infof("Response fm long greet: %v", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	log.Info("Calling greet server bi-di streaming call...")
	reqs := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani",
				LastName:  "Kras",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani2",
				LastName:  "Kras2",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani3",
				LastName:  "Kras3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Dani4",
				LastName:  "Kras4",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	waitc := make(chan struct{})

	// send messages
	go func() {
		for _, req := range reqs {
			log.Infof("Sending req: %v", req)
			stream.Send(req)
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()

	// receive messages
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			log.Infof("Received message from server: %v\n", resp)
		}
		close(waitc)
	}()

	<-waitc
}
