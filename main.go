package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/xans-me/authopia/app"
	"github.com/xans-me/authopia/core/proto"
	"github.com/xans-me/authopia/src/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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

	listerner := app.InjectListener()
	// Inject RPC
	usersRpc, err := users.InjectRPC()
	if err != nil {
		panic(err.Error())
	}
	usersRpc.RegisterRPC(listerner)
	log.Info("Registering Handler User")

	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	mux := runtime.NewServeMux(
		// convert header in response(going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(isHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, writer, request, &newError)
		}))

	log.Info("Registering Handler User")
	proto.RegisterUserServiceHandlerFromEndpoint(
		context.Background(),
		mux,
		"localhost:8081",
		[]grpc.DialOption{grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials())})
	log.Info("Registered Handler User")

	log.Info("Creating a Normal HTTP Server")
	// Creating a normal HTTP server
	server := http.Server{
		Handler: withLogger(mux),
	}
	log.Info("Created a Normal HTTP Server")

	// start server
	err = server.Serve(listerner)
	if err != nil {
		log.Fatal(err)
	}
}

func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func isHeaderAllowed(s string) (string, bool) {
	// check if allowedHeaders contain the header
	if _, isAllowed := allowedHeaders[s]; isAllowed {
		// send uppercase header
		return strings.ToUpper(s), true
	}
	// if not in allowed header, don't send the header
	return s, false
}
