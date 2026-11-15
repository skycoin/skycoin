// Package api src/api/sweep.go c3-coin-api
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

// Sweeper is the node's sweep of exposed addresses. *sweep.Service implements it.
type Sweeper interface {
	SetRule(r sweep.Rule) (sweep.Rule, error)
	Rule(wltID string) (sweep.Rule, bool)
	Status(wltID string) sweep.Status
	Unlock(wltID string, password []byte, ttl time.Duration) error
	Lock(wltID string)
}

// SweepStatus is sweep.Status with the waiting coins as a decimal string.
type SweepStatus struct {
	sweep.Status
	WaitingCoins string `json:"waiting_coins"`
}

// SweepResponse is the sweep rule of a wallet, nil when it has none, and its status.
type SweepResponse struct {
	Rule   *sweep.Rule `json:"rule"`
	Status SweepStatus `json:"status"`
}

// SweepUnlockRequest holds the password of an encrypted wallet in the node so it
// can sign sweeps, for TTLSeconds or, when zero, until the node stops.
type SweepUnlockRequest struct {
	ID         string `json:"id"`
	Password   string `json:"password"`
	TTLSeconds uint64 `json:"ttl_seconds"`
}

// SweepLockRequest makes the node forget the password of a wallet.
type SweepLockRequest struct {
	ID string `json:"id"`
}

func newSweepResponse(s Sweeper, wltID string) (SweepResponse, error) {
	var resp SweepResponse
	if r, ok := s.Rule(wltID); ok {
		resp.Rule = &r
	}
	resp.Status.Status = s.Status(wltID)
	coins, err := droplet.ToString(resp.Status.Status.WaitingCoins)
	if err != nil {
		return resp, err
	}
	resp.Status.WaitingCoins = coins
	return resp, nil
}

func writeSweepError(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	switch {
	case errors.Is(err, wallet.ErrWalletNotExist):
		code = http.StatusNotFound
	case errors.Is(err, wallet.ErrWalletAPIDisabled):
		code = http.StatusForbidden
	case errors.Is(err, wallet.ErrInvalidPassword):
		code = http.StatusUnauthorized
	}
	writeHTTPResponse(w, NewHTTPErrorResponse(code, err.Error()))
}

func writeSweepResponse(w http.ResponseWriter, s Sweeper, wltID string) {
	resp, err := newSweepResponse(s, wltID)
	if err != nil {
		writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	writeHTTPResponse(w, HTTPResponse{Data: resp})
}

// walletSweepHandler gets or sets the sweep rule of a wallet
// URI: /api/v2/wallet/sweep
// Method: GET
// Args:
//
//	id: wallet id [required]
//
// Method: POST
// Body: a sweep.Rule
func walletSweepHandler(s Sweeper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusForbidden, "sweeping needs the wallet API"))
			return
		}
		switch r.Method {
		case http.MethodGet:
			id := r.FormValue("id")
			if id == "" {
				writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, "id is required"))
				return
			}
			writeSweepResponse(w, s, id)
		case http.MethodPost:
			var rule sweep.Rule
			if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
				writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, err.Error()))
				return
			}
			if _, err := s.SetRule(rule); err != nil {
				writeSweepError(w, err)
				return
			}
			writeSweepResponse(w, s, rule.Wallet)
		default:
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""))
		}
	}
}

// walletSweepUnlockHandler holds a wallet password in memory for sweeps
// URI: /api/v2/wallet/sweep/unlock
// Method: POST
// Body: SweepUnlockRequest
func walletSweepUnlockHandler(s Sweeper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusForbidden, "sweeping needs the wallet API"))
			return
		}
		if r.Method != http.MethodPost {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""))
			return
		}
		var req SweepUnlockRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		password := []byte(req.Password)
		req.Password = ""
		defer func() {
			for i := range password {
				password[i] = 0
			}
		}()
		if req.ID == "" {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, "id is required"))
			return
		}
		ttl := time.Duration(req.TTLSeconds) * time.Second //nolint:gosec // seconds from the user
		if err := s.Unlock(req.ID, password, ttl); err != nil {
			writeSweepError(w, err)
			return
		}
		writeSweepResponse(w, s, req.ID)
	}
}

// walletSweepLockHandler forgets a wallet password held for sweeps
// URI: /api/v2/wallet/sweep/lock
// Method: POST
// Body: SweepLockRequest
func walletSweepLockHandler(s Sweeper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s == nil {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusForbidden, "sweeping needs the wallet API"))
			return
		}
		if r.Method != http.MethodPost {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusMethodNotAllowed, ""))
			return
		}
		var req SweepLockRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if req.ID == "" {
			writeHTTPResponse(w, NewHTTPErrorResponse(http.StatusBadRequest, "id is required"))
			return
		}
		s.Lock(req.ID)
		writeSweepResponse(w, s, req.ID)
	}
}

// WalletSweep makes a request to GET /api/v2/wallet/sweep
func (c *Client) WalletSweep(id string) (*SweepResponse, error) {
	v := url.Values{}
	v.Add("id", id)
	var rsp SweepResponse
	ok, err := c.GetV2("/api/v2/wallet/sweep?"+v.Encode(), &rsp)
	if !ok {
		return nil, err
	}
	return &rsp, err
}

// SetWalletSweep makes a request to POST /api/v2/wallet/sweep
func (c *Client) SetWalletSweep(rule sweep.Rule) (*SweepResponse, error) {
	var rsp SweepResponse
	ok, err := c.PostJSONV2("/api/v2/wallet/sweep", rule, &rsp)
	if !ok {
		return nil, err
	}
	return &rsp, err
}

// WalletSweepUnlock makes a request to POST /api/v2/wallet/sweep/unlock
func (c *Client) WalletSweepUnlock(req SweepUnlockRequest) (*SweepResponse, error) {
	var rsp SweepResponse
	ok, err := c.PostJSONV2("/api/v2/wallet/sweep/unlock", req, &rsp)
	if !ok {
		return nil, err
	}
	return &rsp, err
}

// WalletSweepLock makes a request to POST /api/v2/wallet/sweep/lock
func (c *Client) WalletSweepLock(id string) (*SweepResponse, error) {
	var rsp SweepResponse
	ok, err := c.PostJSONV2("/api/v2/wallet/sweep/lock", SweepLockRequest{ID: id}, &rsp)
	if !ok {
		return nil, err
	}
	return &rsp, err
}
