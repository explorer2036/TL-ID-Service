package test

import (
	"TL-ID-Service/proto/id"
	"context"
	"google.golang.org/grpc"
	"testing"
	"time"
	"fmt"
)

func TestGenerate32Bit(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:15002", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := id.NewServiceClient(conn)

	request := id.Generate32BitRequest{}
	reply, err := c.Generate32Bit(context.Background(), &request)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Println(reply)
}
