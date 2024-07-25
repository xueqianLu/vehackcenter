package server

import (
	"context"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"strings"
)

type centerService struct {
	node *Node
	pb.UnimplementedCenterServiceServer
}

func (s *centerService) SubBroadcastTask(in *pb.SubBroadcastTaskRequest, stream pb.CenterService_SubBroadcastTaskServer) error {
	myself := in.Proposer
	ch := make(chan BroadcastEvent)
	sub := s.node.SubscribeBroadcastTask(ch)
	defer sub.Unsubscribe()

	run := true
	for run {
		select {
		case <-stream.Context().Done():
			run = false

		case event := <-ch:
			block := event.Block
			if block.Proposer.Proposer == myself {
				if err := stream.Send(block); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *centerService) SubscribeBlock(in *pb.SubscribeBlockRequest, stream pb.CenterService_SubscribeBlockServer) error {
	myself := strings.ToLower(in.Proposer)
	ch := make(chan NewMinedBlockEvent)
	sub := s.node.SubscribeNewMinedBlock(ch)
	defer sub.Unsubscribe()

	run := true
	for run {
		select {
		case <-stream.Context().Done():
			run = false

		case event := <-ch:
			block := event.Block
			if strings.ToLower(block.Proposer.Proposer) != myself {
				if err := stream.Send(block); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *centerService) SubmitBlock(ctx context.Context, in *pb.Block) (*pb.SubmitBlockResponse, error) {
	s.node.CommitBlock(in)
	return &pb.SubmitBlockResponse{
		Hash: in.Hash,
	}, nil
}

// newCenterServiceServer creates a new CenterServiceServer.
func newCenterServiceServer(node *Node) pb.CenterServiceServer {
	return &centerService{node: node}
}
