package main

import (
	"context"
	"fmt"
	"io"
	"math"
	"net"
	"os"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	calcpb "grpc-course/calc/calc_proto"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type server struct{}

func (*server) CalculateSum(ctx context.Context, req *calcpb.CalculateSumRequest) (*calcpb.CalculateSumResponse, error) {
	log.Infof("Received unary call to calculate sum with request: %v", req)
	return &calcpb.CalculateSumResponse{
		Result: req.X + req.Y,
	}, nil
}

func (*server) PrimeDecompose(req *calcpb.PrimeDecomposeRequest, stream calcpb.Calculator_PrimeDecomposeServer) error {
	log.Infof("Received server stream call to decompose to prime numbers with request: %v", req)

	num := req.Number
	var d int64 = 2
	for num > 1 {
		if num%d == 0 {
			stream.Send(&calcpb.PrimeDecomposeResponse{
				Number: d,
			})
			num = num / d
		} else {
			d = d + 1
		}
	}

	return nil
}

func (*server) CalculateAverage(stream calcpb.Calculator_CalculateAverageServer) error {
	sum := 0.0
	count := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calcpb.CalculateAverageResponse{
				Average: (sum / count),
			})
		}
		if err != nil {
			log.Fatalf("error reading from stream: %v", err)
		}

		sum += req.Number
		count++
	}
}

func (*server) FindMax(stream calcpb.Calculator_FindMaxServer) error {
	log.Info("Processing request for bi dir streaming")
	max := int64(math.MinInt64)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Received error while reading request: %v", err)
			return err
		}
		if n := req.Number; n > max {
			max = n
			err := stream.Send(&calcpb.FindMaxResponse{
				Number: max,
			})
			if err != nil {
				log.Fatalf("Error while sending to stream: %v", err)
				return err
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number %v", number),
		)
	}
	return &calcpb.SquareRootResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	log.Info("Setting up server...")

	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("error listening: %v", err)
	}

	s := grpc.NewServer()

	calcpb.RegisterCalculatorServer(s, &server{})

	// Register reflection service
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("error serving: %v", err)
	}
}
