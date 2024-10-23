package rtc_test

import (
	"context"
	"fmt"
	"gvowr/rtc"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	r, err := rtc.NewClientPeer()
	assert.Nil(t, err)
	if err != nil {
		panic(err)
	}
	err = r.Join("8a296fa2-8f63-11ef-9ef7-88aedd2ba29c")
	assert.Nil(t, err)
	if err != nil {
		panic(err)
	}
	<-r.Cont.Done()
}

func TestCancel(t *testing.T) {
	c, _ := context.WithCancel(context.Background())
	// can()
	<-c.Done()
}

func TestGet(t *testing.T) {
	req := fmt.Sprintf(`{"nodeid":"%v","roomid":"%v"}`, "asdf", "asdf")
	res, _ := http.Post("http://localhost:5050/video/join", "application/json", strings.NewReader(req))
	data, _ := io.ReadAll(res.Body)
	fmt.Println(string(data))

}
