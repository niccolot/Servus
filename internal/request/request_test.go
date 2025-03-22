package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestLineParser(t *testing.T) {
	// test: good GET request line
	req := "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	r, err := RequestFromReader(strings.NewReader(req))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	
	// test: good GET request line with path
	req = "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	r, err = RequestFromReader(strings.NewReader(req))
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// test: invalid number of parts in request line
	req = "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"
	_, err = RequestFromReader(strings.NewReader(req))
	require.Error(t, err)
}
