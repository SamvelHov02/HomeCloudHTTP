package httphelper

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const filePath = "/home/samo/dev/HomeCloud/server/"

type Request struct {
	Method   string
	Resource string
	Headers  Header
	Data     Body
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
	case 204:
		return "No Content"
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
	case 409:
		return "Conflict"
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

type Body struct {
	Data string `json:"data"`
}

func ProcessRequest(conn net.Conn) []byte {
	request := ReadRequest(conn)
	var resp []byte

	switch strings.ToLower(request.Method) {
	case "get":
		data, status, respHeader := ReadGetMethod(request)
		resp = WriteResponse(data, status, respHeader)
	case "post":
		data, status, respHeader := ReadPostMethod(request)
		resp = WriteResponse(data, status, respHeader)
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

	var message string
	var Headers Header
	var Method, Resource string
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		// First line isn't headers
		if len(message) == 0 {
			message = message + line
			lineParts := strings.SplitN(line, " ", 3)
			Method, Resource = lineParts[0], lineParts[1]
			continue
		}

		h := strings.SplitN(line, ":", 2)
		key, val := strings.TrimSpace(h[0]), strings.TrimSpace(h[1])
		Headers.Add(key, val)

		// Header section ended
		if line == "\r\n" {
			break
		}
		message = message + line
	}

	dataLenStr, _ := Headers.Get("Content-Length")
	dataLen, err := strconv.Atoi(dataLenStr[0])

	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, dataLen)

	_, err = io.ReadFull(reader, buf)

	if err != nil {
		log.Fatal(err)
	}

	req := Request{
		Method:   Method,
		Resource: Resource,
		Headers:  Headers,
		Data:     Body{Data: string(buf)},
	}

	return req
}

// Writes Request for some Resource
// Client side code
func WriteRequest(method string, location string, header Header, body Body) []byte {
	var data []byte
	switch strings.ToLower(method) {
	case "get":
		data = WriteGetRequest(location, header)
	case "post":
		data = WritePostRequest(location, header, body)
	}
	return data
}

// Reads the response received from the Server
// Client side code
func ReadResponse(response []byte) (Body, Header, Status) {
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

	var data Body
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
func ReadGetMethod(request Request) ([]byte, Status, Header) {
	var Status Status
	var ResponseHeader Header
	completePath := filePath + request.Resource[1:]
	fmt.Println("Requested file path:", request.Resource)
	fileData, err := os.ReadFile(completePath)
	resp := Body{}

	if errors.Is(err, os.ErrNotExist) {
		Status.Code = 404
	}

	for _, key := range request.Headers.Keys() {
		switch key {
		case "Host":
			continue
		case "Accept":
			val, _ := request.Headers.Get(key)
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
	fmt.Println("Writing GET request for location:", location)
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
func ReadPostMethod(request Request) ([]byte, Status, Header) {
	var Status Status
	var RespHeader Header
	completePath := filePath + request.Resource[1:]

	fmt.Println("POST Request for file path:", request.Resource)

	for _, key := range request.Headers.Keys() {
		switch key {
		case "Host":
			continue
		case "Content-Type":
			if h, _ := request.Headers.Get(key); h[0] != "application/json" {
				Status.Code = 400
			}
		case "Content-Length":
			val, _ := request.Headers.Get(key)
			lenInt, _ := strconv.Atoi(val[0])
			// Body length should be greater than 0
			if lenInt <= 0 {
				Status.Code = 400
			}
		}
	}

	// If something wrong in the headers no need to proceed with body
	if Status.Code == 0 {
		fmt.Println("Gets to create file")
		fmt.Println("Complete Path:", completePath)
		if _, err := os.Stat(completePath); os.IsNotExist(err) {
			if !Info.IsDir() {
				file, err := os.Create(completePath)

				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				file.WriteString(request.Data.Data)
				Status.Code = 204
			} else {
				err := os.Mkdir(completePath, 0755)
			}
		} else {
			// File already exists
			Status.Code = 409

		}
	}

	RespHeader.Add("Content-Length", "0")
	RespHeader.Add("Content-Type", "application/json")

	return []byte{}, Status, RespHeader
}

// Writes Post Requests to send to server
// Client side code
func WritePostRequest(location string, header Header, body Body) []byte {
	fmt.Println("Writing POST request for location:", location)
	dataRaw, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}

	data := []byte("POST " + location + " HTTP/1.1\r\n")

	for _, key := range header.Keys() {
		data = append(data, []byte(key+":"+header[key][0]+"\r\n")...)
	}
	data = append(data, []byte("\r\n")...)
	data = append(data, dataRaw...)
	return data
}

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
