package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "example/gen/greet/v1"        // generated by protoc-gen-go
	"example/gen/greet/v1/greetv1connect" // generated by protoc-gen-connect-go
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())

	if err := ctx.Err(); err != nil {
		return nil, err // automatically coded correctly
	}

	if err := validateGreetRequest(req.Msg); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	greeting, err := doGreetWork(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: greeting,
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func validateGreetRequest(msg *greetv1.GreetRequest) error {
	if msg.Name == "" {
		return errors.New("ぬるぽ")
	}
	return nil
}

func doGreetWork(ctx context.Context, msg *greetv1.GreetRequest) (string, error) {
	if msg.Name == "test" {
		return "", errors.New("test is error")
	}
	return fmt.Sprintf("Hello, %s!", msg.Name), nil
}

func main() {
	// NOTE: muxを入れ子にすることでpathPrefixをつけることができるが、gRPCには対応しなくなる
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(&GreetServer{})
	mux.Handle(path, handler)

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
