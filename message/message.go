package message

import (
	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	return &Message{Body: "Hello From the Server!"}, nil
}