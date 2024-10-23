package api

import (
	"fmt"

	"github.com/samber/lo"
)

type Recommender interface {
	Recommend(roomid string) (*videoNode, error)
}

type MinConnectRecommender struct {
	Server      *GvowrServer
	maxPriority int
}

func (t *MinConnectRecommender) calcPriority(a *videoNode) int {
	return a.nodeDepth * a.peerCount
}
func (t *MinConnectRecommender) Recommend(roomid string) (*videoNode, error) {
	room, ok := t.Server.rooms.M[roomid]
	if !ok {
		var node *videoNode
		return node, fmt.Errorf("not found room")
	}
	nodes := lo.Filter(room.nodes, func(node *videoNode, idx int) bool {
		return node.connectable && t.calcPriority(node) <= t.maxPriority
	})
	if len(nodes) == 0 {
		sourceNode, _ := lo.Find(room.nodes, func(node *videoNode) bool { return node.nodeDepth == 0 })
		return sourceNode, nil
	}
	maxNode := lo.MaxBy(nodes, func(a *videoNode, max *videoNode) bool {
		return t.calcPriority(a) > t.calcPriority(max)
	})
	return maxNode, nil
}
