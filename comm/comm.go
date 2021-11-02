package chat

import (
	"golang.org/x/net/context"
)

type Server struct{
}

func (s *Server) FunTest(ctx context.Context, in *RequestTest) (*ResponseTest, error) {
	return &ResponseTest{Body: "Hello From the Server!"}, nil
}