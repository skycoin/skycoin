package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api/webrpc"
)

func TestWebRPC(t *testing.T) {
	type args struct {
		httpMethod string
		req        webrpc.Request
	}

	cases := []struct {
		name   string
		status int
		args   args
		want   webrpc.Response
	}{
		{
			name:   "http GET",
			status: http.StatusOK,
			args: args{
				httpMethod: http.MethodGet,
				req:        webrpc.Request{},
			},
			want: webrpc.Response{
				Jsonrpc: "2.0",
				Error: &webrpc.RPCError{
					Code:    webrpc.ErrCodeInvalidRequest,
					Message: webrpc.ErrMsgNotPost,
				},
			},
		},

		{
			name:   "invalid jsonrpc",
			status: http.StatusOK,
			args: args{
				httpMethod: http.MethodPost,
				req: webrpc.Request{
					ID:      "1",
					Jsonrpc: "1.0",
					Method:  "get_status",
				},
			},
			want: webrpc.MakeErrorResponse(webrpc.ErrCodeInvalidParams, webrpc.ErrMsgInvalidJsonrpc),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d, err := json.Marshal(tc.args.req)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.args.httpMethod, "/api/v1/webrpc", bytes.NewBuffer(d))
			require.NoError(t, err)

			csrfStore := &CSRFStore{
				Enabled: false,
			}

			gateway := NewGatewayerMock()
			handler := newServerMux(muxConfig{
				host:            configuredHost,
				appLoc:          ".",
				enableJSON20RPC: true,
			}, gateway, csrfStore, nil)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tc.status, rr.Code)

			if rr.Code == http.StatusOK {
				var res webrpc.Response
				err = json.NewDecoder(rr.Body).Decode(&res)
				require.NoError(t, err)
				require.Equal(t, res, tc.want)
			}
		})
	}
}
