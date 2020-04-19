package main

import (
	"fmt"
	"reflect"
	greetpb "grpc-course/greet/greet_pb"
)

func main() {
	r := reflect.TypeOf(greetpb.GreetServiceServer.Greet)
	if r.Out(0).Kind() == reflect.Ptr {
		fmt.Printf("R: %v", r.Out(0).Kind())
	}
}
