// Package electrum implements a JSON-RPC client for Electrum protocol servers (ElectrumX, Fulcrum).
package electrum

import (
	"bufio"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
)

const (
	// DefaultTimeout is the default timeout for RPC requests
	DefaultTimeout = 30 * time.Second
	// ProtocolVersion is the Electrum protocol version we support
	ProtocolVersion = "1.4"
)

// Client is a JSON-RPC client for Electrum protocol servers
type Client struct {
	conn    net.Conn
	scanner *bufio.Scanner
	mu      sync.Mutex
	nextID  atomic.Int64
	timeout time.Duration
}

type jsonRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Result  json.RawMessage `json:"result"`
	// Error is parsed leniently: Electrum servers encode the JSON-RPC error
	// field inconsistently — null/absent/empty on success, a {code,message}
	// object, or (e.g. blockstream's electrs on listunspent/history) a bare
	// string. A fixed *struct type made the whole response fail to unmarshal.
	Error json.RawMessage `json:"error"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// rpcErrorString returns a human error from the JSON-RPC error field, or "" when
// there is no error (null/absent/empty string). Accepts both the {code,message}
// object form and the bare-string form servers use interchangeably.
func rpcErrorString(raw json.RawMessage) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s == "null" || s == `""` {
		return ""
	}
	var obj rpcError
	if err := json.Unmarshal(raw, &obj); err == nil && (obj.Code != 0 || obj.Message != "") {
		return fmt.Sprintf("RPC error %d: %s", obj.Code, obj.Message)
	}
	var str string
	if err := json.Unmarshal(raw, &str); err == nil {
		if strings.TrimSpace(str) == "" {
			return ""
		}
		return "RPC error: " + str
	}
	return "RPC error: " + s
}

// DialFunc dials a raw stream to a "host:port" address. Supplied to
// NewClientWithDialer to route the Electrum connection through a tunnel/proxy
// (e.g. skywire's skysocks-client over the mesh, or a dmsg stream) instead of
// the stdlib clearnet dialer. For ssl:// URLs the TLS handshake still runs
// client-side on top of the returned conn, so the tunnel can't read the stream.
type DialFunc func(network, addr string) (net.Conn, error)

// NewClient creates a new Electrum client connected to the given server URL,
// dialing over the clearnet with the standard library.
// URL format: "ssl://host:port" for TLS or "tcp://host:port" for plain TCP.
func NewClient(serverURL string, timeout time.Duration) (*Client, error) {
	return NewClientWithDialer(serverURL, timeout, nil)
}

// NewClientWithDialer is NewClient with a caller-supplied raw dialer. When
// dial is nil it behaves exactly like NewClient (net.DialTimeout). This lets a
// caller carry the Electrum connection over an arbitrary transport — the whole
// point being that a skywire visor can reach an ssl:// Electrum server through
// skysocks/dmsg while the TLS handshake is still performed end-to-end here.
func NewClientWithDialer(serverURL string, timeout time.Duration, dial DialFunc) (*Client, error) {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	if dial == nil {
		dial = func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		}
	}

	var conn net.Conn
	var err error

	if strings.HasPrefix(serverURL, "ssl://") {
		addr := strings.TrimPrefix(serverURL, "ssl://")
		host, _, herr := net.SplitHostPort(addr)
		if herr != nil {
			return nil, fmt.Errorf("invalid server address %q: %w", addr, herr)
		}
		raw, derr := dial("tcp", addr)
		if derr != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", serverURL, derr)
		}
		tlsConn := tls.Client(raw, &tls.Config{
			ServerName: host,
			MinVersion: tls.VersionTLS12,
			// Electrum servers conventionally present SELF-SIGNED certificates —
			// the Electrum protocol uses TLS for transport encryption, not CA
			// authentication (the reference Electrum client accepts self-signed
			// certs / pins them out of band). CA verification here rejects nearly
			// every public electrum server ("x509: certificate is not valid for
			// any names"), so skip it; ServerName is still sent as SNI.
			InsecureSkipVerify: true, //nolint:gosec // electrum uses self-signed certs by design
		})
		if herr := tlsConn.Handshake(); herr != nil {
			raw.Close() //nolint:errcheck,gosec
			return nil, fmt.Errorf("tls handshake with %s: %w", serverURL, herr)
		}
		conn = tlsConn
	} else if strings.HasPrefix(serverURL, "tcp://") {
		addr := strings.TrimPrefix(serverURL, "tcp://")
		conn, err = dial("tcp", addr)
	} else {
		return nil, fmt.Errorf("unsupported protocol in URL %q (use ssl:// or tcp://)", serverURL)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", serverURL, err)
	}

	client := &Client{
		conn:    conn,
		scanner: bufio.NewScanner(conn),
		timeout: timeout,
	}

	// Increase scanner buffer for large responses
	client.scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	// Perform handshake
	if err := client.handshake(); err != nil {
		conn.Close() //nolint:errcheck,gosec
		return nil, fmt.Errorf("handshake failed: %w", err)
	}

	return client, nil
}

// Close closes the connection
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) call(method string, params []any, result any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	id := c.nextID.Add(1)
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	data = append(data, '\n')

	if err := c.conn.SetDeadline(time.Now().Add(c.timeout)); err != nil {
		return fmt.Errorf("set deadline: %w", err)
	}

	if _, err := c.conn.Write(data); err != nil {
		return fmt.Errorf("write request: %w", err)
	}

	if !c.scanner.Scan() {
		if err := c.scanner.Err(); err != nil {
			return fmt.Errorf("read response: %w", err)
		}
		return fmt.Errorf("connection closed")
	}

	var resp jsonRPCResponse
	if err := json.Unmarshal(c.scanner.Bytes(), &resp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if es := rpcErrorString(resp.Error); es != "" {
		return fmt.Errorf("%s", es)
	}

	if result != nil {
		if err := json.Unmarshal(resp.Result, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}

func (c *Client) handshake() error {
	var result []string
	if err := c.call("server.version", []any{"skycoin-web", ProtocolVersion}, &result); err != nil {
		return err
	}
	if len(result) < 2 {
		return fmt.Errorf("unexpected server.version result: %v", result)
	}
	return nil
}

// BalanceResult is the result of a get_balance call
type BalanceResult struct {
	Confirmed   int64 `json:"confirmed"`
	Unconfirmed int64 `json:"unconfirmed"`
}

// GetBalance returns the balance for an address
func (c *Client) GetBalance(scriptHash string) (*BalanceResult, error) {
	var result BalanceResult
	if err := c.call("blockchain.scripthash.get_balance", []any{scriptHash}, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// UnspentResult is a single UTXO from listunspent
type UnspentResult struct {
	TxHash string `json:"tx_hash"`
	TxPos  uint32 `json:"tx_pos"`
	Height int64  `json:"height"`
	Value  int64  `json:"value"`
}

// ListUnspent returns unspent outputs for a script hash
func (c *Client) ListUnspent(scriptHash string) ([]UnspentResult, error) {
	var result []UnspentResult
	if err := c.call("blockchain.scripthash.listunspent", []any{scriptHash}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// HistoryResult is a single entry from get_history
type HistoryResult struct {
	TxHash string `json:"tx_hash"`
	Height int64  `json:"height"`
	Fee    int64  `json:"fee,omitempty"` // only for mempool txs
}

// GetHistory returns the transaction history for a script hash
func (c *Client) GetHistory(scriptHash string) ([]HistoryResult, error) {
	var result []HistoryResult
	if err := c.call("blockchain.scripthash.get_history", []any{scriptHash}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// BroadcastTransaction broadcasts a raw hex transaction and returns the txid
func (c *Client) BroadcastTransaction(rawTx string) (string, error) {
	var txid string
	if err := c.call("blockchain.transaction.broadcast", []any{rawTx}, &txid); err != nil {
		return "", err
	}
	return txid, nil
}

// EstimateFee returns the estimated fee in BTC per kilobyte for the given number of confirmation blocks
func (c *Client) EstimateFee(numBlocks int) (float64, error) {
	var fee float64
	if err := c.call("blockchain.estimatefee", []any{numBlocks}, &fee); err != nil {
		return 0, err
	}
	return fee, nil
}

// HeaderTip is the current best block header reported by
// blockchain.headers.subscribe.
type HeaderTip struct {
	Height int    `json:"height"`
	Hex    string `json:"hex"`
}

// HeadersSubscribe returns the server's current best block header. It doubles as
// a lightweight liveness probe: a successful call proves the Electrum connection
// works and reports the chain tip height. The subscription side effect (future
// header notifications) is harmless — this synchronous client just ignores them.
func (c *Client) HeadersSubscribe() (*HeaderTip, error) {
	var tip HeaderTip
	if err := c.call("blockchain.headers.subscribe", []any{}, &tip); err != nil {
		return nil, err
	}
	return &tip, nil
}

// GetTransaction returns the raw transaction hex for a given txid
func (c *Client) GetTransaction(txid string) (string, error) {
	var rawTx string
	if err := c.call("blockchain.transaction.get", []any{txid}, &rawTx); err != nil {
		return "", err
	}
	return rawTx, nil
}

// GetTransactionVerbose returns the decoded transaction for a given txid
func (c *Client) GetTransactionVerbose(txid string) (json.RawMessage, error) {
	var result json.RawMessage
	if err := c.call("blockchain.transaction.get", []any{txid, true}, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// AddressToScriptHash converts a Bitcoin address to the Electrum script hash format.
// Supports both P2PKH (1...) and P2WPKH (bc1q...) addresses.
// The script hash is the reversed SHA256 of the output script.
func AddressToScriptHash(address string) (string, error) {
	var script []byte

	if strings.HasPrefix(strings.ToLower(address), "bc1") || strings.HasPrefix(strings.ToLower(address), "tb1") {
		// Native segwit (bech32) P2WPKH address
		segAddr, err := cipher.DecodeBech32BitcoinAddress(address)
		if err != nil {
			return "", fmt.Errorf("invalid bech32 bitcoin address %q: %w", address, err)
		}
		// P2WPKH output script: OP_0 <20-byte witness program>
		script = make([]byte, 22)
		script[0] = 0x00 // OP_0 (witness version 0)
		script[1] = 0x14 // Push 20 bytes
		copy(script[2:22], segAddr.Key[:])
	} else {
		// Legacy P2PKH address
		addr, err := cipher.DecodeBase58BitcoinAddress(address)
		if err != nil {
			return "", fmt.Errorf("invalid bitcoin address %q: %w", address, err)
		}
		// P2PKH output script: OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
		script = make([]byte, 25)
		script[0] = 0x76 // OP_DUP
		script[1] = 0xa9 // OP_HASH160
		script[2] = 0x14 // Push 20 bytes
		copy(script[3:23], addr.Key[:])
		script[23] = 0x88 // OP_EQUALVERIFY
		script[24] = 0xac // OP_CHECKSIG
	}

	// SHA256 of the script
	hash := sha256.Sum256(script)

	// Reverse the hash
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}

	return hex.EncodeToString(hash[:]), nil
}

// AddressToScriptHashFromRaw converts a raw 20-byte pubkey hash to Electrum script hash format
// using P2PKH output script.
func AddressToScriptHashFromRaw(pubKeyHash []byte) (string, error) {
	if len(pubKeyHash) != 20 {
		return "", fmt.Errorf("pubkey hash must be 20 bytes, got %d", len(pubKeyHash))
	}

	// Build P2PKH output script
	script := make([]byte, 25)
	script[0] = 0x76
	script[1] = 0xa9
	script[2] = 0x14
	copy(script[3:23], pubKeyHash)
	script[23] = 0x88
	script[24] = 0xac

	hash := sha256.Sum256(script)
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}

	return hex.EncodeToString(hash[:]), nil
}

// SegwitAddressToScriptHashFromRaw converts a raw 20-byte witness program to Electrum script hash format
// using P2WPKH output script.
func SegwitAddressToScriptHashFromRaw(witnessProgram []byte) (string, error) {
	if len(witnessProgram) != 20 {
		return "", fmt.Errorf("witness program must be 20 bytes, got %d", len(witnessProgram))
	}

	// Build P2WPKH output script: OP_0 <20-byte witness program>
	script := make([]byte, 22)
	script[0] = 0x00 // OP_0
	script[1] = 0x14 // Push 20 bytes
	copy(script[2:22], witnessProgram)

	hash := sha256.Sum256(script)
	for i, j := 0, len(hash)-1; i < j; i, j = i+1, j-1 {
		hash[i], hash[j] = hash[j], hash[i]
	}

	return hex.EncodeToString(hash[:]), nil
}
