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

func getFilter(self string) func(node string) bool {
	return func(node string) bool {
		return strings.ToLower(node) == strings.ToLower(self)
	}
}

func (s *centerService) RegisterNode(ctx context.Context, info *pb.NodeRegisterInfo) (*pb.NodeRegisterResponse, error) {
	s.node.AddRegister(info.Node)
	nodes := s.node.GetAllRegisters(getFilter(info.Node))
	return &pb.NodeRegisterResponse{
		Nodes: nodes,
	}, nil
}

func (s *centerService) FetchNode(ctx context.Context, request *pb.FetchNodeRequest) (*pb.FetchNodeResponse, error) {
	nodes := s.node.GetAllRegisters(getFilter(request.Self))
	return &pb.FetchNodeResponse{
		Nodes: nodes,
	}, nil
}

func (s *centerService) UpdateHack(ctx context.Context, in *pb.UpdateHackRequest) (*pb.Empty, error) {
	s.node.UpdateHack(int(in.Begin), int(in.End))
	return &pb.Empty{}, nil
}

func (s *centerService) SubBroadcastTask(in *pb.SubBroadcastTaskRequest, stream pb.CenterService_SubBroadcastTaskServer) error {

	//myself := in.Proposer
	//ch := make(chan BroadcastEvent, 100)
	//sub := s.node.SubscribeBroadcastTask(ch)
	//defer sub.Unsubscribe()
	//
	//run := true
	//for run {
	//	select {
	//	case <-stream.Context().Done():
	//		run = false
	//
	//	case event := <-ch:
	//		block := event.Block
	//		if block.Proposer.Proposer == myself {
	//			if err := stream.Send(block); err != nil {
	//				return err
	//			}
	//		}
	//	}
	//}
	return nil
}

func (s *centerService) SubscribeBlock(in *pb.SubscribeBlockRequest, stream pb.CenterService_SubscribeBlockServer) error {
	myself := strings.ToLower(in.Proposer)
	ch := make(chan NewBlockEvent, 100)
	sub := s.node.SubscribeNewBlock(ch)
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

func (s *centerService) SubscribeMinedBlock(in *pb.SubscribeBlockRequest, stream pb.CenterService_SubscribeMinedBlockServer) error {
	myself := strings.ToLower(in.Proposer)
	ch := make(chan NewMinedBlockEvent, 100)
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

func (s *centerService) BroadcastBlock(ctx context.Context, block *pb.Block) (*pb.SubmitBlockResponse, error) {
	s.node.BroadcastBlock(block)
	return &pb.SubmitBlockResponse{
		Hash: block.Hash,
	}, nil
}

func (s *centerService) SubmitBlock(ctx context.Context, in *pb.Block) (*pb.SubmitBlockResponse, error) {
	s.node.CommitBlock(in)
	return &pb.SubmitBlockResponse{
		Hash: in.Hash,
	}, nil
}

func (s *centerService) Vote(ctx context.Context, in *pb.VoteRequest) (*pb.VoteResponse, error) {
	value := s.node.conf.Vote
	if in.Block < int64(s.node.conf.BeginToHack) || in.Block > int64(s.node.conf.EndToHack) {
		value = 1
	}
	return &pb.VoteResponse{
		Vote: int32(value),
	}, nil
}

// newCenterServiceServer creates a new CenterServiceServer.
func newCenterServiceServer(node *Node) pb.CenterServiceServer {
	return &centerService{node: node}
}
