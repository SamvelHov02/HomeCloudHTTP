package httphelper

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var filePath = "/home/samo/dev/HomeCloud/server/Vault/"

type Request struct {
	Method   string
	Resource string
	Headers  Header
}

type Status struct {
	Code int
}

func (s Status) Text() string {
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

func (h Header) Get(key string) ([]string, bool) {
	val, ok := h[key]
	return val, ok
}

func (h *Header) Add(key string, value string) {
	if *h == nil {
		*h = make(map[string][]string)
	}
	(*h)[key] = append((*h)[key], value)
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

func ProcessRequest(conn net.Conn) []byte {
	request := ReadRequest(conn)
	var resp []byte

	switch strings.ToLower(request.Method) {
	case "get":
		data, status, respHeader := ReadGetMethod(request.Resource, request.Headers)
		resp = WriteResponse(data, status, respHeader)
	case "post":
		readPostMethod(request.Resource, request.Headers)
	case "put":
		readPutMethod(request.Resource, request.Headers)
	case "delete":
		readDeleteMethod(request.Resource, request.Headers)
	}
	return resp
}

// Read and extracts meta information from request
// Server side code

func ReadRequest(conn net.Conn) Request {
	reader := bufio.NewReader(conn)

	var headers Header
	headersFinished := false
	var message string
	lenInt := 0
	var body []byte
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		if len(message) == 0 {
			message = message + line
			continue
		}

		if line != "\r\n" && !headersFinished {
			message = message + line
			h := strings.SplitN(line, ":", 2)
			key, val := strings.TrimSpace(h[0]), strings.TrimSpace(h[1])
			headers.Add(key, val)
		} else {
			headersFinished = true
			message = message + line
			if len, ok := headers.Get("Content-Length"); ok {
				lenInt, _ = strconv.Atoi(len[0])
				body = make([]byte, lenInt)
				_, err = io.ReadFull(reader, body)
				if err != nil {
					log.Fatal(err)
				}
				message = message + string(body)
			}
			break
		}
	}
	meta := strings.SplitN(message, " ", 3)

	request := Request{
		Method:   meta[0],
		Resource: meta[1],
		Headers:  headers,
	}

	return request
}

// Writes Request for some Resource
// Client side code
func WriteRequest(method string, location string, header Header) []byte {
	var data []byte
	switch strings.ToLower(method) {
	case "get":
		data = WriteGetRequest(location, header)
	}
	return data
}

// Reads the response received from the Server
// Client side code
func ReadResponse(response []byte) (ResponseBody, Header, Status) {
	// Header and body seperated by \n so a \n\n sequence indicates end of headers
	strResponse := string(response)
	parts := strings.Split(strResponse, "\r\n\r\n")
	headersField, body := parts[0], []byte(parts[1])
	var status Status
	var headers Header
	for i, line := range strings.Split(headersField, "\r\n") {
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
func WriteResponse(data []byte, Status Status, headers Header) []byte {
	resp := []byte("HTTP/1.1 " + strconv.Itoa(Status.Code) + " " + Status.Text() + "\r\n")

	for _, key := range headers.Keys() {
		switch key {
		case "Content-Type":
			line := []byte("Content-Type:" + headers[key][0] + "\r\n")
			resp = append(resp, line...)
		case "Content-Length":
			line := []byte("Content-Length:" + headers[key][0] + "\r\n")
			resp = append(resp, line...)
		}
	}
	line := []byte("Server:HomeCloud/0.0.1\r\n\r\n")
	resp = append(resp, line...)
	resp = append(resp, data...)
	return resp
}

// Reads Get requests from the Clients,
// Server side code
func ReadGetMethod(uri string, headers Header) ([]byte, Status, Header) {
	var Status Status
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
			val, _ := headers.Get(key)
			if val[0] == "application/json" {
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
	data := []byte("GET " + "/" + location + " HTTP/1.1\r\n")

	for _, key := range header.Keys() {
		if len(header[key]) == 1 {
			data = append(data, []byte(key+":"+header[key][0]+"\r\n")...)
		}
	}
	data = append(data, []byte("\r\n")...)
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
