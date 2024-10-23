package api

import (
	"fmt"
)

func (t *GvowrServer) setVideoApi() {
	videoGroup := t.Group("/video")
	videoGroup.POST("/new", jsonFunc(t.api_videoNew))
	videoGroup.POST("/join", jsonFunc(t.api_videoJoin))
}

func (t *GvowrServer) api_videoNew(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	nodeId := reqMap["nodeid"].(string)

	room := newVideoRoom()
	node := newVideoNode()
	node.nodeid = nodeId
	node.connectable = true
	node.nodeDepth = 0
	room.nodes = append(room.nodes, node)
	roomId := makeUuid()

	t.rooms.M[roomId] = room
	resMap["roomid"] = roomId
	fmt.Printf("t.videos: %v\n", t.rooms)
	return nil
}

func (t *GvowrServer) api_videoJoin(reqMap map[string]any, resMap map[string]any) error {
	t.rooms.Lock.Lock()
	defer t.rooms.Lock.Unlock()

	roomid := reqMap["roomid"].(string)
	nodeid := reqMap["nodeid"].(string)
	room, ok := t.rooms.M[roomid]
	if !ok {
		return fmt.Errorf(roomid, " not found")
	}
	node := newVideoNode()
	node.nodeid = nodeid
	recommendNode, err := t.recommender.Recommend(roomid)
	if err != nil {
		return err
	}
	room.nodes = append(room.nodes, node)
	resMap["nodeid"] = recommendNode.nodeid
	return nil
}

type videoNode struct {
	nodeid      string
	peer        *videoNode
	peerCount   int
	nodeDepth   int
	connectable bool
}
type videoRoom struct {
	nodes []*videoNode
}

func newVideoNode() *videoNode {
	return &videoNode{peerCount: 0}
}
func newVideoRoom() *videoRoom {
	return &videoRoom{
		nodes: make([]*videoNode, 0),
	}
}
