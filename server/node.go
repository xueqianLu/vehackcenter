package server

import (
	"github.com/xueqianLu/vehackcenter/config"
	"github.com/xueqianLu/vehackcenter/event"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"log"
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

	mux       sync.Mutex
	registers map[string]string

	conf config.Config
}

func NewNode(conf config.Config) *Node {
	n := &Node{
		conf:      conf,
		registers: make(map[string]string),
	}
	maxMsgSize := 20 * 1024 * 1024
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
	// 1. send block to all subscribed hackers.
	n.minedBlockFeed.Send(NewMinedBlockEvent{Block: block})
	T := int64(10)

	go func(b *pb.Block) {
		/*
		 * begin = t - T * index
		 * end = begin + (2n - 1) * T
		 */
		blockTime := int64(b.Timestamp)
		begin := blockTime - T*int64(b.Proposer.Index)
		end := begin + (2*int64(n.conf.HackerCount)-1)*T
		targetBlockTime := end
		next := targetBlockTime // 出块者会提前5秒开始出块，在这里提前5秒广播
		if b.Height < int64(n.conf.BeginToHack) {
			next = 1
		} else {
			time.Sleep(time.Duration(next-time.Now().Unix()) * time.Second)
		}

		log.Printf("time to release block %d, by proposer %s\n", b.Height, b.Proposer.Proposer)

		// 3. then send BroadcastTask to proposer
		n.broadcastTaskFeed.Send(BroadcastEvent{Block: b})
	}(block)

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

	log.Println("server start at ", n.conf.Url)

	if err := n.apiServer.Serve(lis); err != nil {
		log.Printf("grpc serve err: %v", err)
	}
}

func (n *Node) StopServer() {
	n.apiServer.Stop()
}
