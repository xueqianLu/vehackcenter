package server

import (
	"github.com/xueqianLu/vehackcenter/config"
	"github.com/xueqianLu/vehackcenter/event"
	pb "github.com/xueqianLu/vehackcenter/hackcenter"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Node struct {
	broadcastTaskFeed event.Feed
	minedBlockFeed    event.Feed
	scope             event.SubscriptionScope
	apiServer         *grpc.Server

	conf config.Config
}

func NewNode(conf config.Config) *Node {
	n := &Node{
		conf: conf,
	}
	maxMsgSize := 20 * 1024 * 1024
	// create grpc server
	n.apiServer = grpc.NewServer(grpc.MaxSendMsgSize(maxMsgSize), grpc.MaxRecvMsgSize(maxMsgSize))
	return n
}

func (n *Node) CommitBlock(block *pb.Block) {
	// 1. send block to all subscribed hackers.
	n.minedBlockFeed.Send(NewMinedBlockEvent{Block: block})

	go func(b *pb.Block) {
		nextNBlockTime := func(curBlockTime int64, nextCount int64) int64 {
			return curBlockTime + int64(10)*nextCount
		}
		blockTime := int64(b.Timestamp)
		targetBlockTime := nextNBlockTime(blockTime, int64(n.conf.HackerCount-int(b.Proposer.Index)))
		next := targetBlockTime - 5 // 出块者会提前5秒开始出块，在这里提前5秒广播
		if b.Height < int64(n.conf.BeginToHack) {
			next = 1
		} else {
			time.Sleep(time.Duration(next-time.Now().Unix()) * time.Second)
		}

		for {
			// broad cast before next block 500ms.
			now := time.Now().UnixMilli()
			if now >= (int64(next)*1000 + 300) {
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		log.Printf("time to release block %d, by proposer %s\n", b.Height, b.Proposer.Proposer)

		// 3. then send BroadcastTask to proposer
		n.broadcastTaskFeed.Send(BroadcastEvent{Block: b})
	}(block)

}

func (n *Node) SubscribeNewMinedBlock(ch chan<- NewMinedBlockEvent) event.Subscription {
	return n.scope.Track(n.minedBlockFeed.Subscribe(ch))
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
