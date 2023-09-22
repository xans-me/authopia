package main

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"github.com/soheilhy/cmux"
	"github.com/xans-me/authopia/app"
	"github.com/xans-me/authopia/core/proto"
	http2 "github.com/xans-me/authopia/helpers/http"
	"github.com/xans-me/authopia/src/users"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	// initialize config
	config, err := app.InitializeAppConfig()
	if err != nil {
		log.Panic(err)
	}

	// Inject RPC
	usersRpc, err := users.InjectRPC()
	if err != nil {
		log.Fatal(err)
	}

	log.Info("############################")
	log.Info("AUTHOPIA")
	log.Info("You app running on ", config.App.Environment, " mode")
	log.Info("############################")

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Creating mux for gRPC gateway. This will multiplex or route request different gRPC service
	grpcMux := runtime.NewServeMux(jsonOption,
		// convert header in response(going from gateway) from metadata received.
		runtime.WithOutgoingHeaderMatcher(http2.IsHeaderAllowed),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			// send all the headers received from the client
			md := metadata.Pairs("auth", header)
			return md
		}),
		runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaller runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {
			//creating a new HTTTPStatusError with a custom status, and passing error
			newError := runtime.HTTPStatusError{
				HTTPStatus: 400,
				Err:        err,
			}
			// using default handler to do the rest of heavy lifting of marshaling error and adding headers
			runtime.DefaultHTTPErrorHandler(ctx, mux, marshaller, writer, request, &newError)
		}))

	err = proto.RegisterUserServiceHandlerServer(context.Background(), grpcMux, usersRpc)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// for handling swagger
	// statikFS, err := fs.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	// mux.Handle("/swagger/", swaggerHandler)

	// Creating a normal HTTP server
	httpServer := http.Server{
		Handler: http2.WithLogger(mux),
	}

	// get listener
	listener := app.InjectListener()

	// Setup multiplexer connection
	m := cmux.New(listener)
	grpcListener := m.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	//groups of goroutines working on subtask registry of multi connection
	g := new(errgroup.Group)
	g.Go(func() error { return usersRpc.RegisterRPC(grpcListener) })
	g.Go(func() error { return httpServer.Serve(httpListener) })
	g.Go(func() error { return m.Serve() })

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
