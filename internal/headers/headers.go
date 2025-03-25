package headers

import (
	"fmt"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {	
	headersString := string(data)

	// all headers are parsed
	if strings.HasPrefix(headersString, "\r\n\r\n") {
		return 4, true, nil
	}	

	crlfindex := strings.Index(headersString, "\r\n")
	
	// need more data to find a complete header
	if crlfindex == -1 {
		return 0, false, nil
	}

	headerLine := headersString[:crlfindex]
	headerLine = strings.TrimSpace(headerLine)
	
	if strings.Contains(headerLine, "::") {
		return 0, false, fmt.Errorf("invalid header format: double colon found")
	}

	colonIndex := strings.Index(headerLine, ":")
	if colonIndex == -1 {
		return 0, false, fmt.Errorf("invalid header format: no colon found")
	}

	key := headerLine[:colonIndex]

	if !isValidHeaderFieldName(key) {
		return 0, false, fmt.Errorf("invalid character in field name")
	}

	key = strings.ToLower(key)
	if strings.TrimSpace(key) != key {
		return 0, false, fmt.Errorf("invalid header format: space before colon found")
	}

	value := strings.TrimSpace(headerLine[colonIndex + 1:])
	value = strings.Trim(value, ",")  

	prevVal, ok := h[key]

	// if this field name is already in the headers map, append the value
	// to the preexisting one separated by a comma and whitespace
	if ok {
		h[key] = prevVal + ", " + value
	} else {
		h[key] = value
	}
	
	// +2 for \r\n
	return crlfindex + 2, false, nil
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
