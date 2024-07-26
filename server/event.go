package server

import pb "github.com/xueqianLu/vehackcenter/hackcenter"

type NewMinedBlockEvent struct{ Block *pb.Block }

type NewBlockEvent struct{ Block *pb.Block }

type BroadcastEvent struct{ Block *pb.Block }
