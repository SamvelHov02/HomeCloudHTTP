package httphelper

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var filePath = "/Users/samvelhovhannisyan/Documents/dev/Personal/HomeCloud/Vault"

type HTTPStatus struct {
	Code int
}

func (s HTTPStatus) Text() string {
	switch s.Code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 301:
		return "Moved Permanently"
	case 302:
		return "Found"
	case 400:
		return "Bad Request"
	case 401:
		return "Unauthorized"
	case 404:
		return "Not Found"
	default:
		return ""
	}
}

type Header map[string][]string

func (h Header) Get(key string) []string {
	return h[key]
}

func (h Header) Add(key string, value string) {
	if h == nil {
		h = make(map[string][]string)
	}
	h[key] = append(h[key], value)
}

func (h Header) Keys() []string {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	return keys
}

type ResponseBody struct {
	Data string `json:"data"`
}

func ProcessRequest(conn net.Conn, request string) []byte {
	method, uri, headers := ReadRequest(request)
	var resp []byte

	switch method {
	case "GET":
		data, status, respHeader := readGetMethod(uri, headers)
		resp = WriteResponse(conn, data, status, respHeader)
	case "POST":
		readPostMethod(uri, headers)
	case "PUT":
		readPutMethod(uri, headers)
	case "DELETE":
		readDeleteMethod(uri, headers)
	}
	return resp
}

// Read and extracts meta information from request
// Server side code
func ReadRequest(request string) (string, string, Header) {
	lines := strings.Split(request, "\n")
	requestLine := strings.Split(lines[0], " ")
	method, uri := requestLine[0], requestLine[1]
	var headers Header
	for _, line := range lines[1:] {
		h := strings.SplitN(line, ":", 2)
		key, val := strings.TrimSpace(h[0]), strings.TrimSpace(h[1])
		for _, v := range strings.Split(val, ",") {
			headers.Add(key, v)
		}
	}
	return method, uri, headers
}

// Writes Request for some Resource
// Client side code
func WriteRequest(method int, location string, header Header) []byte {
	var data []byte
	switch method {
	case 1:
		data = WriteGetRequest(location, header)
	}
	return data
}

// Reads the response received from the Server
// Client side code
func ReadResponse(response []byte) (ResponseBody, Header, HTTPStatus) {
	// Header and body seperated by \n so a \n\n sequence indicates end of headers
	strResponse := string(response)
	parts := strings.Split(strResponse, "\n\n")
	headersField, body := parts[0], []byte(parts[1])
	var status HTTPStatus
	var headers Header
	for i, line := range strings.Split(headersField, "\n") {
		if i == 0 {
			components := strings.Split(line, " ")
			if components[0] != "HTTP/1.1" {
				log.Fatal("Bad Response 1")
			}

			if code, err := strconv.Atoi(components[1]); err == nil {
				status.Code = code
			} else {
				log.Fatal("Bad Response 2")
			}

			if codeNum := status.Text(); codeNum != components[2] {
				log.Fatal("Bad Response 3")
			}
		} else {
			h := strings.SplitN(line, ":", 2)
			key, val := strings.TrimSpace(h[0]), strings.TrimSpace(h[1])
			for _, v := range strings.Split(val, ",") {
				headers.Add(key, v)
			}
		}
	}

	var data ResponseBody
	err := json.Unmarshal(body, &data)

	if err != nil {
		log.Fatal(err)
	}

	return data, headers, status
}

// Writes Response to the request
// Server side code
func WriteResponse(conn net.Conn, data []byte, Status HTTPStatus, headers Header) []byte {
	resp := []byte("HTTP/1.1 " + strconv.Itoa(Status.Code) + Status.Text() + "\n")

	for _, key := range headers.Keys() {
		switch key {
		case "Content-Type":
			line := []byte("Content-Type:" + headers[key][0] + "\n")
			resp = append(resp, line...)
		case "Content-Length":
			line := []byte("Content-Length:" + headers[key][0] + "\n")
			resp = append(resp, line...)
		}
		line := []byte("Server:HomeCloud/0.0.1\n")
		resp = append(resp, line...)
	}
	resp = append(resp, []byte("\n")...)
	resp = append(resp, data...)
	return resp
}

// Reads Get requests from the Clients,
// Server side code
func readGetMethod(uri string, headers Header) ([]byte, HTTPStatus, Header) {
	var Status HTTPStatus
	var ResponseHeader Header
	completePath := filePath + uri
	fileData, err := os.ReadFile(completePath)
	resp := ResponseBody{}

	if err == os.ErrNotExist {
		Status.Code = 401
	}

	for _, key := range headers.Keys() {
		switch key {
		case "Host":
			continue
		case "Accept":
			if headers.Get(key)[0] == "application/json" {
				resp.Data = string(fileData)
				ResponseHeader.Add("Content-Type", "application/json")
				Status.Code = 200
			}
			// case "Authorization":
			// continue
		}
	}
	dataJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	ResponseHeader.Add("Content-Length", strconv.Itoa(len(dataJSON)))

	return dataJSON, Status, ResponseHeader
}

// Writes Get Requests to send to server
// Client side code
func WriteGetRequest(location string, header Header) []byte {
	data := []byte("GET " + location + " HTTP/1.1\n")

	for _, key := range header.Keys() {
		if len(header[key]) == 1 {
			data = append(data, []byte(key+":"+header[key][0]+"\n")...)
		}
	}
	return data
}

// Reads Post requests from the Clients,
// Server side code
func readPostMethod(uri string, headers Header) {}

// Writes Post Requests to send to server
// Client side code
func WritePostRequest() {}

// Reads Put requests from the Clients,
// Server side code
func readPutMethod(uri string, headers Header) {}

// Writes Put Requests to send to server
// Client side code
func WritePutRequest() {}

// Reads Delete requests from the Clients,
// Server side code
func readDeleteMethod(uri string, headers Header) {}

// Writes Delete Requests to send to server
// Client side code
func WriteDeleteRequest() {}
