package api

import context "context"

type Server struct{}

func (s *Server) GetStatus(context.Context, *StatusRequest) (*StatusReply, error) {
	return &StatusReply{Status: "ok"}, nil
}
