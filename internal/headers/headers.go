package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func (h* Headers) Parse(data []byte) (n int, done bool, err error) {	
	headersString := string(data)
	if strings.Contains(headersString, "::") {
		return 0, false, fmt.Errorf("invalid header: double colon")
	}

	crlfIndex := strings.Index(headersString, "\r\n")
	
	// not enough data to have a full header
	if crlfIndex == -1 {
		return 0, false, nil
	}

	// headers are done and empty line reached, consume the crlf
	if crlfIndex == 0 {
		return 2, true, nil
	}

	parts := strings.SplitN(headersString[:crlfIndex], ":", 2)
	key := strings.ToLower(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := strings.TrimSpace(parts[1])

	// remove commas from data to avoid possibly malformed values
	value = strings.Trim(value, ",")
	key = strings.TrimSpace(key)

	if !isValidHeaderFieldName(key) {
		return 0, false, fmt.Errorf("invalid header token: %s", key)
	}

	h.Add(key, value)
	
	// +2 for \r\n
	return crlfIndex + 2, false, nil
}

func (h* Headers) Add(key, value string) {
	/*
	*@brief: add a (key,value) pair and append the additional value
	* if key is already present
	*/
	prevVal, ok := h.Get(key)
	if ok {
		(*h)[key] = prevVal + ", " + value
	} else {
		(*h)[key] = value
	}
}

func (h* Headers) AddOverride(key, value string) {
	/*
	*@brief: add a (key, value) pair overriding an eventual
	* preexisting value
	*/
	(*h)[key] = value
}

func (h *Headers) Get(key string) (string, bool) {
	val, ok := (*h)[strings.ToLower(key)]

	return val, ok
}

func GetDefaultHeaders(contentLen int) Headers {
	headers := Headers{}
	headers.Add("Content-Length", fmt.Sprint(contentLen))
	headers.Add("Connection", "close")
	headers.Add("Content-Type", "text/plain")

	return headers
}

func isValidHeaderFieldName(s string) bool {
	/*
	* field names must contain only:
	* uppercase or lowercase letters, 0-9 digits
	* ! # $ % & ' * + - . ^ _ ` | ~ special characters
	*/
	for _, ch := range s {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') ||
			(ch >= '0' && ch <= '9') || strings.ContainsRune("-!#$%&'*+.^_`|~", ch)) {

			return false
		}
	}

	return true
}
