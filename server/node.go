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
		begin := blockTime - T*int64(block.Proposer.Index)
		log.WithFields(log.Fields{
			"begin":    begin,
			"block":    block.Height,
			"proposer": block.Proposer.Proposer,
			"index":    block.Proposer.Index,
		}).Info("CommitBlock receive")

		var newList []*pb.Block = nil
		n.mux.Lock()
		if _, exist := n.hackedBlockList[begin]; !exist {
			n.hackedBlockList[begin] = make([]*pb.Block, 0)
			newList = n.hackedBlockList[begin]
		}
		n.hackedBlockList[begin] = append(n.hackedBlockList[begin], block)
		n.mux.Unlock()
		log.WithField("pending", len(n.hackedBlockList[begin])).Info("CommitBlock pendin block count")

		go func(begin int64, list []*pb.Block) {
			/*
			 * begin = t - T * index
			 * end = begin + (2n - 1) * T
			 */
			if list == nil {
				return
			}

			end := begin + (2*int64(n.conf.HackerCount)-1)*T
			targetBlockTime := end
			next := targetBlockTime // 出块者会提前5秒开始出块，在这里提前5秒广播
			log.WithFields(log.Fields{
				"begin": begin,
				"wait":  next - time.Now().Unix(),
			}).Info("CommitBlock wait to broadcast hacked block")
			time.Sleep(time.Duration(next-time.Now().Unix()) * time.Second)

			mlist := n.hackedBlockList[begin]
			for _, b := range mlist {
				log.WithFields(log.Fields{
					"hacked-block": b.Height,
					"proposer":     b.Proposer.Proposer,
				}).Info("CommitBlock time to release hacked block")
				n.BroadcastBlock(b)
			}
			log.WithField("begin", begin).Info("CommitBlock broadcast hacked block finished")
		}(begin, newList)
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
