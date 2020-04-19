package main

import (
	"context"
	"fmt"
	calcpb "grpc-course/calc/calc_proto"
	"io"
	"math"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("Setting up client...")

	cc, err := grpc.Dial("127.0.0.1:80", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error dialing: %s", err)
	}
	defer cc.Close()

	c := calcpb.NewCalculatorClient(cc)

	doUnaryCall(c)
	doServerStreamCall(c)
	doClientStreamCall(c)
	doBiDirStreamCall(c)
	doErrorUnary(c)
}

func doErrorUnary(c calcpb.CalculatorClient) {
	log.Info("Calling calc unary error with a positive number...")
	resp, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: 225})
	if err != nil {
		responseError, ok := status.FromError(err)
		if ok {
			// actual error form grpc (user error)
			fmt.Println(responseError.Message())
			fmt.Println(responseError.Code())
			if responseError.Code() == codes.InvalidArgument {
				fmt.Println("We have probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big error calling squareroot: %v", err)
		}
	}
	log.Infof("Response from calc unary: %s", resp)

	log.Info("Calling calc unary error with a negative number...")
	respTwo, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: -12})
	if err != nil {
		responseError, ok := status.FromError(err)
		if ok {
			// actual error form grpc (user error)
			fmt.Println(responseError.Message())
			fmt.Println(responseError.Code())
			if responseError.Code() == codes.InvalidArgument {
				fmt.Println("We have probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big error calling squareroot: %v", err)
		}
	}
	log.Infof("Response from calc unary: %s", respTwo)
}

func doUnaryCall(c calcpb.CalculatorClient) {
	log.Info("Calling calc unary...")

	req := &calcpb.CalculateSumRequest{
		X: 21,
		Y: 11,
	}

	resp, err := c.CalculateSum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling calc unary service: %v", err)
	}

	log.Infof("Response from calc unary: %s", resp)
}

func doServerStreamCall(c calcpb.CalculatorClient) {
	log.Info("Calling calc server stream...")

	req := &calcpb.PrimeDecomposeRequest{
		Number: 465723,
	}

	respStream, err := c.PrimeDecompose(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling calc unary service: %v", err)
	}
	for {
		resp, err := respStream.Recv()
		if err == io.EOF {
			log.Info("Reached end of stream for result of call to server streaming...")
			break
		}
		if err != nil {
			log.Fatalf("Error while calling calc server stream service: %v", err)
		}
		log.Infof("Response from calc server stream call: %s", resp)
	}
}

func doClientStreamCall(c calcpb.CalculatorClient) {
	log.Info("Calling calc client stream...")
	stream, err := c.CalculateAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling calc server client streaming: %v", err)
	}
	for i := 0; i < 115; i++ {
		ff := float64(i) / math.Pi * math.E
		err := stream.Send(&calcpb.CalculateAverageRequest{
			Number: ff,
		})
		if err != nil {
			log.Fatalf("Smth very bad happened while calling client stream: %v", err)
		}
		log.Infof("Calling client stream with: %v\n", ff)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream: %v", err)
	}
	log.Infof("Grande resultante: %v\n", res)
}

func doBiDirStreamCall(c calcpb.CalculatorClient) {
	log.Info("Calling calc bi dir stream...")
	stream, err := c.FindMax(context.Background())
	if err != nil {
		log.Fatalf("Error while getting stream: %v", err)
	}

	waitc := make(chan struct{})
	reqs := []int64{12, -2, 1233, 1, 7, -86, 346, 92, 477, 485644, 67564468, 576, 213446}
	go func() {
		for _, n := range reqs {
			err := stream.Send(&calcpb.FindMaxRequest{
				Number: n,
			})
			log.Infof("Sending bi dir message to server: %v", n)
			if err != nil {
				log.Fatalf("Error sending to server: %v", err)
			}
			time.Sleep(time.Second)
		}
		log.Infof("Closing stream")
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading message: %v", err)
				break
			}
			log.Infof("Received message from server: %v", resp)
		}
		close(waitc)
	}()

	<-waitc
}
