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
	broadcastTaskFeed    event.Feed
	minedBlockFeed       event.Feed
	newBlockFeed         event.Feed
	newExternalBlockFeed event.Feed
	scope                event.SubscriptionScope
	apiServer            *grpc.Server
	hackedBlockList      map[int64][]*pb.Block
	pendingBlockChan     chan *pb.Block

	quit chan struct{}

	mux       sync.Mutex
	registers map[string]string

	conf config.Config
}

func NewNode(conf config.Config) *Node {
	n := &Node{
		conf:             conf,
		registers:        make(map[string]string),
		hackedBlockList:  make(map[int64][]*pb.Block),
		pendingBlockChan: make(chan *pb.Block, 10000),
		quit:             make(chan struct{}),
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

func (n *Node) broadCastPending() {
	externalBlockCh := make(chan NewBlockEvent, 10000)
	sub := n.SubscribeNewExternalBlock(externalBlockCh)
	defer sub.Unsubscribe()

	var realBroadcast = func(duration int64, blks []*pb.Block) {
		time.Sleep(time.Duration(duration) * time.Second)
		for _, blk := range blks {
			n.BroadcastBlock(blk, true)
			log.WithFields(log.Fields{
				"height":   blk.Height,
				"hash":     blk.Hash,
				"proposer": blk.Proposer.Proposer,
			}).Info("Broadcast pending hacker block")
		}
	}

	var newBlock *pb.Block
	T := 10
	var firstPending *pb.Block

	var toBroadcast []*pb.Block
	for {
		select {
		case <-n.quit:
			return
		case ev := <-externalBlockCh:
			// got a new honest block, clear pending list and broadcast all pending blocks.
			newBlock = ev.Block
			if firstPending == nil {
				log.WithFields(log.Fields{
					"height":   newBlock.Height,
					"proposer": newBlock.Proposer.Proposer,
					"pending":  len(toBroadcast),
					"fist":     "null",
				}).Info("New honest block arrived, no pending block")
				continue
			}
			if newBlock.Height <= firstPending.Height {
				log.WithFields(log.Fields{
					"height":   newBlock.Height,
					"proposer": newBlock.Proposer.Proposer,
					"pending":  len(toBroadcast),
					"fist":     firstPending.Height,
				}).Info("New honest block arrived with little height, ignore it")
				continue
			}
			pendLength := len(toBroadcast)
			if pendLength > 0 {
				// calculate the duration to broadcast pending blocks.
				targetTime := newBlock.Timestamp + int64(pendLength*T)
				end := targetTime - 10 // before 6 seconds of the target block.
				duration := end - time.Now().Unix()
				if duration < 0 {
					duration = 0
				}
				log.WithFields(log.Fields{
					"height":       newBlock.Height,
					"proposer":     newBlock.Proposer.Proposer,
					"pending":      pendLength,
					"firstPending": firstPending.Height,
					"targetTime":   targetTime,
				}).Info("New honest block arrived, broadcast pending blocks")
				realBroadcast(duration, toBroadcast)
				toBroadcast = make([]*pb.Block, 0)
			}

		case newPending := <-n.pendingBlockChan:
			if len(toBroadcast) == 0 {
				firstPending = newPending
			}
			toBroadcast = append(toBroadcast, newPending)
		}

	}
}

func (n *Node) CommitBlock(block *pb.Block) {
	if block.Height < int64(n.conf.BeginToHack) || block.Height > int64(n.conf.EndToHack) {
		// direct broadcast to all nodes.
		n.BroadcastBlock(block, true)
	} else {
		// 1. send block to all subscribed hackers.
		n.minedBlockFeed.Send(NewMinedBlockEvent{Block: block})
		// add to pending list.
		n.pendingBlockChan <- block
	}
}

func (n *Node) BroadcastBlock(block *pb.Block, internal bool) {
	// 1. send block to all subscribed node.
	n.newBlockFeed.Send(NewBlockEvent{Block: block})
	if !internal {
		// 2. send block to all subscribed external node.
		n.newExternalBlockFeed.Send(NewBlockEvent{Block: block})
	}
}

func (n *Node) SubscribeNewExternalBlock(ch chan<- NewBlockEvent) event.Subscription {
	return n.scope.Track(n.newExternalBlockFeed.Subscribe(ch))
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

	go n.broadCastPending()

	if err := n.apiServer.Serve(lis); err != nil {
		log.WithError(err).Error("grpc serve error")
	}
}

func (n *Node) StopServer() {
	n.apiServer.Stop()
	close(n.quit)
}

func (n *Node) UpdateHack(begin int, end int) {
	n.conf.BeginToHack = begin
	n.conf.HackerCount = end
}
