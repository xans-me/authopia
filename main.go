package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/app"
	"github.com/xans-me/authopia/core/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

func main() {
	// initialize config
	config, err := app.InitializeAppConfig()
	if err != nil {
		panic(err.Error())
	}

	log.Info("############################")
	log.Info("AUTHOPIA")
	log.Info("You app running on ", config.App.Environment, " mode")
	log.Info("############################")

	// get listener
	//listener := app.InjectListener()

	//userRPC, err := users.InjectRPC()
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}

	conn, err := grpc.DialContext(
		context.Background(),
		"0.0.0.0"+config.App.Port,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = proto.RegisterUserServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":50052",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway for REST on http://0.0.0.0:50052")
	if err := gwServer.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway server: %v", err)
	}

}
