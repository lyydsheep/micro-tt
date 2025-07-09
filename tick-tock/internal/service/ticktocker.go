package service

import (
	"context"

	pb "tick-tock/api/tick_tock/v1"
)

type TickTockerService struct {
	pb.UnimplementedTickTockerServer
}

func NewTickTockerService() *TickTockerService {
	return &TickTockerService{}
}

func (s *TickTockerService) CreateTaskDefine(ctx context.Context, req *pb.CreateRequest) (*pb.CreateReply, error) {
	return &pb.CreateReply{}, nil
}
func (s *TickTockerService) UpdateTaskDefineStatus(ctx context.Context, req *pb.UpdateRequest) (*pb.UpdateReply, error) {
	return &pb.UpdateReply{}, nil
}
