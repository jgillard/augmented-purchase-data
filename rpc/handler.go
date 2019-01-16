package api

import (
	"context"

	"github.com/jgillard/practising-go-tdd/internal"
)

type Server struct{}

func (s *Server) GetStatus(context.Context, *StatusRequest) (*StatusReply, error) {
	status := internal.GetStatus()
	return &StatusReply{Status: status}, nil
}
