package headers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// test: valid single header
func TestHeaderParser(t *testing.T) {
	os.Stdout.Sync()
	headers := Headers{}
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// test: invalid spacing header
	headers = Headers{}
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test: valid single header with extra withespaces
	headers = Headers{}
	data = []byte("   Host: localhost:42069   \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 29, n)
	assert.False(t, done)

	// test: invalid character in field name
	headers = Headers{}
	data = []byte("H@st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test: valid 2 headers
	headers = Headers{}
	data = []byte("Host: localhost:42069\r\nContent-Type: application-json\r\n\r\n")
	n1, done, err := headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n1) 
	assert.False(t, done)
	n2, done, err := headers.Parse(data[n1:]) // parse second header
	require.NoError(t, err)
	assert.Equal(t, "application-json", headers["content-type"])
	assert.Equal(t, 32, n2) 
	assert.False(t, done)

	// test: single header with multiple values
	headers = Headers{
		"content-type": "application-json",
	}

	data = []byte("Host: localhost:42069\r\nContent-Type: text/html\r\n\r\n")
	n1, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n1) 
	assert.False(t, done)
	n2, done, err = headers.Parse(data[n1:]) // parse second header
	require.NoError(t, err)
	assert.Equal(t, "application-json, text/html", headers["content-type"])
	assert.Equal(t, 25, n2) 
	assert.False(t, done)

	// test: empty field
	headers = Headers{}
	data = []byte("Host:\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "", headers["host"])
	assert.Equal(t, 7, n)
	assert.False(t, done)

	// test: missing colon
	headers = Headers{}
	data = []byte("Host localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test: double colon
	headers = Headers{}
	data = []byte("Host:: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// test: field value with leading comma
	headers = Headers{}
	data = []byte("Host: ,localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 24, n)
	assert.False(t, done)
}
