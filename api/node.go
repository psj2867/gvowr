package api

import (
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

func (t *GvowrServer) setNodeApi() {
	nodeGroup := t.Group("/node/connect")
	nodeGroup.POST("/", jsonFunc(t.api_nodeConnectInfo))
	nodeGroup.POST("/success", jsonFunc(t.api_nodeConnectSuccess))
	nodeGroup.POST("/fail/src", jsonFunc(t.api_nodeSourceDelete))
	nodeGroup.POST("/fail/view", jsonFunc(t.api_nodeViewerDelete))
}

func (t *GvowrServer) api_nodeConnectSuccess(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	roomid := reqMap["roomid"].(string)
	nodeid := reqMap["nodeid"].(string)
	remote := reqMap["remote"].(string)
	room, ok := t.rooms.M[roomid]
	if !ok {
		return fmt.Errorf("room not found")
	}
	node, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == nodeid
	})
	if !ok {
		return fmt.Errorf("node not found")
	}
	remoteNode, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == remote
	})
	if !ok {
		return fmt.Errorf("remote node not found")
	}
	node.peer = remoteNode
	node.connectable = true
	node.nodeDepth = remoteNode.nodeDepth + 1
	remoteNode.peerCount++
	node.peerCount = 1
	return nil
}

func (t *GvowrServer) api_nodeSourceDelete(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	roomid := reqMap["roomid"].(string)
	nodeid := reqMap["nodeid"].(string)
	remote := reqMap["remote"].(string)
	room, ok := t.rooms.M[roomid]
	if !ok {
		return fmt.Errorf("room not found")
	}
	node, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == nodeid
	})
	if !ok {
		return fmt.Errorf("node not found")
	}
	remoteNode, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == remote
	})
	if !ok {
		return fmt.Errorf("remote node not found")
	}
	node.peerCount--
	remoteNode.connectable = false
	return nil
}

func (t *GvowrServer) api_nodeViewerDelete(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	roomid := reqMap["roomid"].(string)
	remote := reqMap["remote"].(string)
	room, ok := t.rooms.M[roomid]
	if !ok {
		return fmt.Errorf("room not found")
	}
	remoteNode, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == remote
	})
	if !ok {
		return fmt.Errorf("remote node not found")
	}
	remoteNode.connectable = false
	return nil
}

func (t *GvowrServer) api_nodeConnectInfo(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	roomid := reqMap["roomid"].(string)
	nodeid := reqMap["nodeid"].(string)
	connectCount := reqMap["count"]
	room, ok := t.rooms.M[roomid]
	if !ok {
		return fmt.Errorf("room not found")
	}
	node, ok := lo.Find(room.nodes, func(node *videoNode) bool {
		return node.nodeid == nodeid
	})
	if !ok {
		return fmt.Errorf("node not found")
	}
	count, ok := parseInt(connectCount)
	if !ok {
		return fmt.Errorf("count is not int")
	}
	node.peerCount = count
	return nil
}

func parseInt(v any) (int, bool) {
	switch t := v.(type) {
	case int, int8, int16, int32, int64:
		a := reflect.ValueOf(t).Int() // a has type int64
		return int(a), true
	case uint, uint8, uint16, uint32, uint64:
		a := reflect.ValueOf(t).Uint() // a has type uint64
		return int(a), true
	}
	return 0, false
}
