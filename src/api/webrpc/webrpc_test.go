package webrpc

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/stretchr/testify/assert"
)

type fakeGateway struct {
}

func (fg fakeGateway) GetLastBlocks(num uint64) ([]coin.Block, error) {
	return []coin.Block{}, nil
}

func setup() (*rpcHandler, func()) {
	c := make(chan struct{})
	f := func() {
		close(c)
	}

	return newRPCHandler(1, 1, &fakeGateway{}, c), f
}

func TestHTTPMethod(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()
	d, err := json.Marshal(Request{})
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("GET", "/webrpc", bytes.NewBuffer(d))
	w := httptest.NewRecorder()
	rpc.Handler(w, r)

	var res Response
	if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, res.Error.Code, errCodeRequirePost)
}

func TestGetStatus(t *testing.T) {
	rpc, teardown := setup()
	defer teardown()

	testDatas := []struct {
		HTTPMethod    string
		Req           *Request
		ExpectErrCode int
	}{
	// {
	// 	"GET",
	// 	nil,
	// 	errCodeRequirePost,
	// },
	}

	res := Response{}
	for _, data := range testDatas {
		d, err := json.Marshal(data.Req)
		if err != nil {
			t.Fatal(err)
		}

		r := httptest.NewRequest(data.HTTPMethod, "/webrpc", bytes.NewBuffer(d))
		w := httptest.NewRecorder()
		rpc.Handler(w, r)
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, data.ExpectErrCode, res.Error.Code)
	}
}
