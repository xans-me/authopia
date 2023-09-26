package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/xans-me/authopia/app"
	"github.com/xans-me/authopia/core/proto"
	http2 "github.com/xans-me/authopia/helpers/http"
	"github.com/xans-me/authopia/helpers/response"
	"github.com/xans-me/authopia/src/users"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// Initialize config
	config, err := app.InitializeAppConfig()
	if err != nil {
		log.Fatal("Failed to initialize app config:", err)
	}

	// Inject RPC
	usersRpc, err := users.InjectRPC()
	if err != nil {
		log.Fatal("Failed to inject RPC:", err)
	}

	log.Info("############################")
	log.Info("AUTHOPIA")
	log.Info("Your app is running in", config.App.Environment, "mode")
	log.Info("############################")

	// JSON Marshalling options
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Create a ServeMux for gRPC Gateway
	grpcMux := runtime.NewServeMux(jsonOption,
		// runtime.WithForwardResponseOption(GRPCGatewayHTTPResponseModifier),
		// Convert headers in response (going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(http2.IsHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// Send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			// Use custom error message if the error is a CustomError
			if customErr, ok := err.(*response.ErrorStruct); ok {
				// Convert the CustomError to JSON
				errorJSON, _ := json.Marshal(customErr)
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusInternalServerError) // You can set the appropriate HTTP status code here
				writer.Write(errorJSON)
				return
			}

			// Use the default handler to handle other errors
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaller, writer, request, err)
		}))
	// Register the gRPC UserService handler
	err = proto.RegisterUserServiceHandlerServer(context.Background(), grpcMux, usersRpc)
	if err != nil {
		log.Fatal("Failed to register UserService handler:", err)
	}

	// Create a normal HTTP server
	httpServer := http.Server{
		Handler: http2.WithLogger(grpcMux),
	}

	// Get the listener
	listener := app.InjectListener()

	// Setup multiplexer connection using cmux
	m := cmux.New(listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	// Groups of goroutines working on subtask registry of multi-connection
	g := new(errgroup.Group)
	g.Go(func() error { return usersRpc.RegisterRPC(grpcListener) })
	g.Go(func() error { return httpServer.Serve(httpListener) })
	g.Go(func() error { return m.Serve() })

	// Wait for all goroutines to finish and handle any errors
	if err := g.Wait(); err != nil {
		log.Fatal("Error occurred:", err)
	}
}
