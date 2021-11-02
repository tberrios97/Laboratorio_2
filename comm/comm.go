package chat

import (
	"golang.org/x/net/context"
)

type Server struct{
}

func (s *Server) funTest(ctx context.Context, in *requestTest) (*responseTest, error) {
	return &responseTest{Body: "Hello From the Server!"}, nil
}