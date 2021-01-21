package main

import (
	"example/handler"
	pb "example/proto"
	"time"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("example"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterExampleHandler(srv.Server(), new(handler.Example))

	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	// let's delay the process for exiting for reasons you'll see below
	time.Sleep(time.Second * 5)
}
