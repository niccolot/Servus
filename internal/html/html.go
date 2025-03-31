package html

import (

	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"golang.org/x/net/html"

	"Servus/internal/headers"
	"Servus/internal/response"
)

func WriteResponse(w *response.Writer, fileName string) error {
	status, body, err := extractTitleAndBodyFromFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to parse html: %v", err)
	}
	code, err := extractCode(status)
	if err != nil {
		fmt.Errorf("failed to parse status code: %v", err)
	}

	w.Response.Code = response.StatusCode(code)
	w.Response.Message = body
	w.Response.Headers = headers.Headers{}
	w.Response.Headers.AddOverride("Content-Type", "text/html")
	w.Response.Headers.AddOverride("Content-Length", fmt.Sprint(len(w.Response.Message)))

	return nil
}

func extractTitleAndBodyFromFile(htmlFile string) (string, []byte, error) {
	file, err := os.Open(htmlFile)
	if err != nil {
		return "", []byte{}, fmt.Errorf("Error opening file:", err)
	}
	defer file.Close()

	// Read file contents
	content, err := io.ReadAll(file)
	if err != nil {
		return "", []byte{}, fmt.Errorf("Error reading file:", err)
		
	}

	title, body, err := extractTitleAndBody(content)
	if err != nil {
		return "", []byte{}, fmt.Errorf("Error parsing HTML:", err)
	}

	return title, body, nil
}

func extractCode(s string) (int, error) {
	re := regexp.MustCompile(`[0-9]+`) 
	numStr := re.FindString(s)
	if numStr == "" {
		return 0, fmt.Errorf("status code not found")
	}

	return strconv.Atoi(numStr)
}

func extractTitleAndBody(htmlStr []byte) (string, []byte, error) {
	doc, err := html.Parse(strings.NewReader(string(htmlStr)))
	if err != nil {
		return "", []byte{}, err
	}

	var title, body string
	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		if n.Type == html.ElementNode && n.Data == "body" {
			body = extractText(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(doc)
	return title, []byte(body), nil
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var text string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c) + " "
	}
	return strings.Join(strings.Fields(text), " ")
}
