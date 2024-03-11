package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	greet "github.com/aminnouri-com/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greet.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, request *greet.GreetRequest) (*greet.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", request)
	firstName := request.GetGreeting().GetFirstName()
	result := "hello " + firstName

	res := &greet.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(request *greet.GreetManytimesRequest, stream greet.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", request)
	firstName := request.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " Number " + strconv.Itoa(i)
		response := &greet.GreetManytimesResponse{
			Result: result,
		}

		stream.Send(response)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greet.GreetService_LongGreetServer) error {
	fmt.Println("LogGreet function was invoked with Streaming request")
	result := ""
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greet.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("%v", err)
		}

		firstName := request.Greeting.GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greet.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone function wass invoked with streaming request")

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client streaming request: %v", err)
			return err
		}
		firstName := request.GetGreeting().GetFirstName()
		restul := "Hello " + firstName
		sendErr := stream.Send(&greet.GreetEveryoneResponse{
			Result: restul,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", err)
		}
	}
}

func main() {
	fmt.Println("Greet Server Listening...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greet.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
