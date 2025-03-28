package request

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequestLineParser(t *testing.T) {
	// test: good GET request line
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}

	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)
	
	// test: good GET request line with path
	reader = &chunkReader{
		data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// test: good GET request line with path, max numBytesPerRead
	reader = &chunkReader{
		data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "GET", r.RequestLine.Method)
	require.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	require.Equal(t, "1.1", r.RequestLine.HttpVersion)

	// test: invalid number of parts in request line
	reader = &chunkReader{
		data: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data) + 1,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// test: invalid method (lowercase)
	reader = &chunkReader{
		data: "get /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data),
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// test: invalid method (non http method)
	reader = &chunkReader{
		data: "gimme /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data),
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// test: invalid method (out of order)
	reader = &chunkReader{
		data: "/coffee GET HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data),
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// test: invalid target
	reader = &chunkReader{
		data: "GET / coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data),
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)

	// test: invalid http version
	reader = &chunkReader{
		data: "GET / coffee HTTP/1.0\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: len(reader.data),
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestHeadersParser(t *testing.T) {
	// test: standard headers
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	require.Equal(t, "localhost:42069", r.Headers["host"])
	require.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	require.Equal(t, "*/*", r.Headers["accept"])

	// test: malformed header
	reader = &chunkReader{
		data: "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	_, err = RequestFromReader(reader)
	require.Error(t, err)
}

func TestBodyParser(t *testing.T) {
//	// test: standard body
//	reader := &chunkReader{
//		data: "POST /submit HTTP/1.1\r\n" +
//			"Host: localhost: 42069\r\n" +
//			"Content-Length: 13\r\n" +
//			"\r\n" +
//			"hello world!\n",
//		numBytesPerRead: 3,
//	}
//	r, err := RequestFromReader(reader)
//	require.NoError(t, err)
//	require.NotNil(t, r)
//	require.Equal(t, "hello world!\n", string(r.Body))
//
//	// test: no body, 'Content-Length: 0'
//	reader = &chunkReader{
//		data: "POST /submit HTTP/1.1\r\n" +
//			"Host: localhost: 42069\r\n" +
//			"Content-Length: 0\r\n" +
//			"\r\n",
//		numBytesPerRead: 3,
//	}
//	r, err = RequestFromReader(reader)
//	require.NoError(t, err)
//	require.NotNil(t, r)
//	require.Equal(t, "", string(r.Body))
//
//	// test: no body, no 'Content-Length'
//	reader = &chunkReader{
//		data: "POST /submit HTTP/1.1\r\n" +
//			"Host: localhost: 42069\r\n" +
//			"\r\n",
//		numBytesPerRead: 3,
//	}
//	r, err = RequestFromReader(reader)
//	require.NoError(t, err)
//	require.NotNil(t, r)
//	require.Equal(t, "", string(r.Body))
//
//	// test: body shorter than 'Content-Length'
//	reader = &chunkReader{
//		data: "POST /submit HTTP/1.1\r\n" +
//			"Host: localhost: 42069\r\n" +
//			"Content-Length: 13\r\n" +
//			"\r\n" +
//			"hello wo\n",
//		numBytesPerRead: 3,
//	}
//	_, err = RequestFromReader(reader)
//	require.Error(t, err)
//
//	// test: body longer than 'Content-Length'
	reader := &chunkReader{
		data: "POST /submit HTTP/1.1\r\n" +
			"Host: localhost: 42069\r\n" +
			"Content-Length: 13\r\n" +
			"\r\n" +
			"hello world!!!!!!\n",
		numBytesPerRead: 1,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
}
