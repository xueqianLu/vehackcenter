package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/xueqianLu/vehackcenter/config"
	"github.com/xueqianLu/vehackcenter/event"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type Node struct {
	broadcastTaskFeed event.Feed
	minedBlockFeed    event.Feed
	newBlockFeed      event.Feed
	scope             event.SubscriptionScope
	apiServer         *grpc.Server
	hackedBlockList   map[int64][]*pb.Block

	mux       sync.Mutex
	registers map[string]string

	conf config.Config
}

func NewNode(conf config.Config) *Node {
	n := &Node{
		conf:            conf,
		registers:       make(map[string]string),
		hackedBlockList: make(map[int64][]*pb.Block),
	}
	maxMsgSize := 100 * 1024 * 1024
	// create grpc server
	n.apiServer = grpc.NewServer(grpc.MaxSendMsgSize(maxMsgSize), grpc.MaxRecvMsgSize(maxMsgSize))
	return n
}

func (n *Node) AddRegister(enode string) {
	n.mux.Lock()
	defer n.mux.Unlock()
	n.registers[enode] = enode
}

func (n *Node) GetAllRegisters(filter func(node string) bool) []string {
	n.mux.Lock()
	defer n.mux.Unlock()
	var registers []string
	for _, v := range n.registers {
		if filter(v) {
			continue
		}
		registers = append(registers, v)
	}
	return registers
}

func (n *Node) CommitBlock(block *pb.Block) {

	T := int64(10)
	if block.Height < int64(n.conf.BeginToHack) || block.Height > int64(n.conf.EndToHack) {
		// direct broadcast to all nodes.
		n.BroadcastBlock(block)
	} else {
		// 1. send block to all subscribed hackers.
		n.minedBlockFeed.Send(NewMinedBlockEvent{Block: block})

		// add to hack block list, and when the time is up, broadcast the block.
		blockTime := int64(block.Timestamp)
		next := blockTime + T
		end := next - 3

		log.WithFields(log.Fields{
			"block":             block.Height,
			"proposer":          block.Proposer.Proposer,
			"index":             block.Proposer.Index,
			"wait to broadcast": end - time.Now().Unix(),
		}).Info("CommitBlock receive")

		go func(duration int64, blk *pb.Block) {
			time.Sleep(time.Duration(duration) * time.Second)
			log.WithFields(log.Fields{
				"hacked-block": blk.Height,
				"proposer":     blk.Proposer.Proposer,
			}).Info("CommitBlock time to release hacked block")
			n.BroadcastBlock(blk)
		}(end-time.Now().Unix(), block)
	}
}

func (n *Node) BroadcastBlock(block *pb.Block) {
	// 1. send block to all subscribed node.
	n.newBlockFeed.Send(NewBlockEvent{Block: block})
}

func (n *Node) SubscribeNewMinedBlock(ch chan<- NewMinedBlockEvent) event.Subscription {
	return n.scope.Track(n.minedBlockFeed.Subscribe(ch))
}

func (n *Node) SubscribeNewBlock(ch chan<- NewBlockEvent) event.Subscription {
	return n.scope.Track(n.newBlockFeed.Subscribe(ch))
}

func (n *Node) SubscribeBroadcastTask(ch chan<- BroadcastEvent) event.Subscription {
	return n.scope.Track(n.broadcastTaskFeed.Subscribe(ch))
}

func (n *Node) RunServer() {
	// listen port
	lis, err := net.Listen("tcp", n.conf.Url)
	if err != nil {
		log.Printf("listen port err: %v", err)
		return
	}

	// register service into grpc server
	pb.RegisterCenterServiceServer(n.apiServer, newCenterServiceServer(n))

	log.WithField("url", n.conf.Url).Info("server start")

	if err := n.apiServer.Serve(lis); err != nil {
		log.WithError(err).Error("grpc serve error")
	}
}

func (n *Node) StopServer() {
	n.apiServer.Stop()
}

func (n *Node) UpdateHack(begin int, end int) {
	n.conf.BeginToHack = begin
	n.conf.HackerCount = end
}
