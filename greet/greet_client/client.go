package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	greet "github.com/aminnouri-com/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client.")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greet.NewGreetServiceClient(cc)

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBiDiStreaming(c)
}

func doUnary(c greet.GreetServiceClient) {
	fmt.Println("Strating to do an Unary RPC")
	request := &greet.GreetRequest{
		Greeting: &greet.Greeting{
			FirstName: "Amin",
			LastName:  "Nouri",
		},
	}

	response, err := c.Greet(context.Background(), request)

	if err != nil {
		log.Fatalf("error while call greet RPC: %v", err)
	}

	log.Printf("Response from greet: %v", response.Result)
}

func doServerStreaming(c greet.GreetServiceClient) {
	fmt.Println("Starting to a server streaming RPC")

	request := &greet.GreetManytimesRequest{
		Greeting: &greet.Greeting{
			FirstName: "Amin",
			LastName:  "Nouri",
		},
	}
	streamResponse, err := c.GreetManyTimes(context.Background(), request)
	if err != nil {
		log.Fatalf("error while call greet manytimes RPC %v\n", err)
	}

	for {
		message, err := streamResponse.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManytimes: %v", message.GetResult())
	}
}

func doClientStreaming(c greet.GreetServiceClient) {
	fmt.Println("Starting a client streaming RPC")

	requests := []*greet.LongGreetRequest{
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Amin",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Azi",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Mani",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Peyman",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Mehrbod",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("errot while calling greet long: %v", err)
	}

	for _, request := range requests {
		fmt.Printf("Sending request: %v\n", request)
		stream.Send(request)
		time.Sleep(1000 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receving response LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", response)
}

func doBiDiStreaming(c greet.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("errot while creating stream: %v", err)
		return
	}

	requests := []*greet.GreetEveryoneRequest{
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Amin",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Azi",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Mani",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Peyman",
			},
		},
		&greet.GreetEveryoneRequest{
			Greeting: &greet.Greeting{
				FirstName: "Mehrbod",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, request := range requests {
			fmt.Printf("Sending message: %v\n", request)
			stream.Send(request)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Errot while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", response.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
