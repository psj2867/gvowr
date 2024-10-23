package rtc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	peer "github.com/muka/peerjs-go"
)

type PeerConnection struct {
	conn *peer.Peer
	id   string
	Cont context.Context
}

func NewClientPeer() (*PeerConnection, error) {
	cont, cancel := context.WithCancel(context.Background())
	options := peer.NewOptions()
	options.Host = "n1.psj2867.com"
	options.Port = 9000
	options.Secure = false
	options.Debug = 3
	uuid, _ := uuid.NewUUID()
	id := uuid.String()
	clinetPeer, err := peer.NewPeer(id, options)
	if err != nil {
		return nil, err
	}
	clinetPeer.On("connection", func(data interface{}) {
		conn := data.(*peer.DataConnection)
		conn.On("data", func(i interface{}) {
			fmt.Println("open")
			go conn.Close()
		})
	})
	clinetPeer.On("call", func(data interface{}) {
		conn := data.(*peer.MediaConnection)
		options := peer.AnswerOption{}
		conn.Answer(nil, &options)
		conn.On("stream", func(i interface{}) {
			fmt.Printf("i: %v\n", i)
			fmt.Println("stream connected")
			cancel()
		})
	})
	return &PeerConnection{conn: clinetPeer, id: id, Cont: cont}, nil
}

func (t *PeerConnection) Join(roomid string) error {
	req := fmt.Sprintf(`{"nodeid":"%v","roomid":"%v"}`, t.id, roomid)
	res, err := http.Post("http://localhost:5050/video/join", "application/json", strings.NewReader(req))
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("stats code %v", res.StatusCode)
	}
	resMap := make(map[string]any)
	dataB, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataB, &resMap)
	if err != nil {
		return err
	}
	nodeid, ok := resMap["nodeid"].(string)
	if !ok {
		return fmt.Errorf("nodeid is invalid")
	}
	options := peer.NewConnectionOptions()
	options.Debug = 3

	dataConn, err := t.conn.Connect(nodeid, options)
	if err != nil {
		return err
	}
	dataConn.On("open", func(i interface{}) {
		err = dataConn.Send([]byte(t.id), false)
	})
	return err
}
