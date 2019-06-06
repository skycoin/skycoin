// +build base extra full

package base

import (
	// "fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const (
	// StatusCode = 4
	// Body 4
	RESPONSE_SIZE = 8
)

func op_http_get(expr *CXExpression, fp int) {
	inp1, out1 := expr.Inputs[0], expr.Outputs[0]
	out1Offset := GetFinalOffset(fp, out1)

	var err error
	var resp *http.Response
	var contents []byte

	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err = netClient.Get(ReadStr(fp, inp1))

	// resp, err = http.Get(ReadStr(mem, fp, inp1))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	contents, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	byts := encoder.Serialize(string(contents))
	length := len(byts)
	heapOffset := AllocateSeq(length + OBJECT_HEADER_SIZE)
	size := encoder.Serialize(int32(len(byts)))

	var header []byte = make([]byte, OBJECT_HEADER_SIZE, OBJECT_HEADER_SIZE)
	for c := 5; c < OBJECT_HEADER_SIZE; c++ {
		header[c] = size[c-5]
	}

	obj := append(header, byts...)

	WriteMemory(heapOffset, obj)

	off := encoder.SerializeAtomic(int32(heapOffset))

	WriteMemory(out1Offset, off)
}

// type Response struct {
//         Status     string // e.g. "200 OK"
//         StatusCode int    // e.g. 200
//         Proto      string // e.g. "HTTP/1.0"
//         ProtoMajor int    // e.g. 1
//         ProtoMinor int    // e.g. 0

//         // Header maps header keys to values. If the response had multiple
//         // headers with the same key, they may be concatenated, with comma
//         // delimiters.  (Section 4.2 of RFC 2616 requires that multiple headers
//         // be semantically equivalent to a comma-delimited sequence.) When
//         // Header values are duplicated by other fields in this struct (e.g.,
//         // ContentLength, TransferEncoding, Trailer), the field values are
//         // authoritative.
//         //
//         // Keys in the map are canonicalized (see CanonicalHeaderKey).
//         Header Header

//         // Body represents the response body.
//         //
//         // The response body is streamed on demand as the Body field
//         // is read. If the network connection fails or the server
//         // terminates the response, Body.Read calls return an error.
//         //
//         // The http Client and Transport guarantee that Body is always
//         // non-nil, even on responses without a body or responses with
//         // a zero-length body. It is the caller's responsibility to
//         // close Body. The default HTTP client's Transport may not
//         // reuse HTTP/1.x "keep-alive" TCP connections if the Body is
//         // not read to completion and closed.
//         //
//         // The Body is automatically dechunked if the server replied
//         // with a "chunked" Transfer-Encoding.
//         Body io.ReadCloser

//         // ContentLength records the length of the associated content. The
//         // value -1 indicates that the length is unknown. Unless Request.Method
//         // is "HEAD", values >= 0 indicate that the given number of bytes may
//         // be read from Body.
//         ContentLength int64

//         // Contains transfer encodings from outer-most to inner-most. Value is
//         // nil, means that "identity" encoding is used.
//         TransferEncoding []string

//         // Close records whether the header directed that the connection be
//         // closed after reading Body. The value is advice for clients: neither
//         // ReadResponse nor Response.Write ever closes a connection.
//         Close bool

//         // Uncompressed reports whether the response was sent compressed but
//         // was decompressed by the http package. When true, reading from
//         // Body yields the uncompressed content instead of the compressed
//         // content actually set from the server, ContentLength is set to -1,
//         // and the "Content-Length" and "Content-Encoding" fields are deleted
//         // from the responseHeader. To get the original response from
//         // the server, set Transport.DisableCompression to true.
//         Uncompressed bool

//         // Trailer maps trailer keys to values in the same
//         // format as Header.
//         //
//         // The Trailer initially contains only nil values, one for
//         // each key specified in the server's "Trailer" header
//         // value. Those values are not added to Header.
//         //
//         // Trailer must not be accessed concurrently with Read calls
//         // on the Body.
//         //
//         // After Body.Read has returned io.EOF, Trailer will contain
//         // any trailer values sent by the server.
//         Trailer Header

//         // Request is the request that was sent to obtain this Response.
//         // Request's Body is nil (having already been consumed).
//         // This is only populated for Client requests.
//         Request *Request

//         // TLS contains information about the TLS connection on which the
//         // response was received. It is nil for unencrypted responses.
//         // The pointer is shared between responses and should not be
//         // modified.
//         TLS *tls.ConnectionState
// }
