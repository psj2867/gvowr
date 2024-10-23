package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gvowr/api"
	"gvowr/util"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var s *api.GvowrServer

func TestMain(m *testing.M) {
	s = api.Server()
	os.Exit(m.Run())
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	res := httptest.NewRecorder()
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)

}

func TestNewVideo(t *testing.T) {
	req := requestClient("POST", "/video/new", map[string]any{
		"nodeid": "asdf",
	})
	res := httptest.NewRecorder()
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
	fmt.Printf("jsonUnMarshalIo(res.Body): %v\n", jsonUnMarshalIo(res.Body))
}

func TestJoinVideo(t *testing.T) {
	req := requestClient("POST", "/video/new", map[string]any{
		"nodeid": "node1",
	})
	res := httptest.NewRecorder()
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
	resMap := jsonUnMarshalIo(res.Body)

	roomid := resMap["roomid"].(string)
	req = requestClient("POST", "/video/join", map[string]any{
		"nodeid": "node2",
		"roomid": roomid,
	})
	res = httptest.NewRecorder()
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code)
	resMap = jsonUnMarshalIo(res.Body)
	fmt.Printf("resMap: %v\n", resMap)
}

func TestJoinVideoFail(t *testing.T) {

	req := requestClient("POST", "/video/join", map[string]any{
		"nodeid": "node2",
		"roomid": "asdf",
	})
	res := httptest.NewRecorder()
	s.ServeHTTP(res, req)
	assert.NotEqual(t, 200, res.Code)

}

func TestSyncMapFail(t *testing.T) {
	a := util.SyncMap[string, int]{}
	_, ok := a.Load("")
	assert.Equal(t, false, ok)
}

func jsonMarshalIo(data map[string]any) io.Reader {
	return bytes.NewReader(jsonMarshal(data))
}
func jsonMarshal(data map[string]any) []byte {
	r, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return r
}
func jsonUnMarshalIo(data *bytes.Buffer) map[string]any {
	l := data.Len()
	b := make([]byte, l)
	_, err := data.Read(b)
	if err != nil {
		panic(err)
	}
	return jsonUnMarshal(b)
}
func jsonUnMarshal(data []byte) map[string]any {
	res := make(map[string]any, 0)
	err := json.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	return res
}

func requestClient(method string, url string, data map[string]any) *http.Request {
	req := httptest.NewRequest(method, url, jsonMarshalIo(data))
	req.Header.Add("Content-Type", "application/json")
	return req
}
