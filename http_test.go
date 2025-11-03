package httphelper

import (
	"fmt"
	"strconv"
	"testing"
)

func TestReadResponse(t *testing.T) {
	response := []byte("HTTP/1.1 200 OK\r\nContent-Type: application/json\r\nContent-Length: 11\r\n\r\n{\"data\":\"Hi\"}")
	actualResponse, _, actualStatus := ReadResponse(response)
	expectedResponse := Body{
		Data: "Hi",
	}
	fmt.Println("Does it even get here")
	//expectedHeaders := Header{
	//"Content-Type":   []string{"application/json"},
	//"Content-Length": []string{"11"},
	//}

	expectedStatus := Status{Code: 200}

	if actualResponse.Data != expectedResponse.Data {
		t.Errorf("Got %s, expected %s", actualResponse.Data, expectedResponse.Data)
	}

	// if !reflect.DeepEqual(actualHeaders, expectedHeaders) {
	// fmt.Println(actualHeaders)
	// fmt.Println("---------------------")
	// fmt.Println(expectedHeaders)
	// t.Errorf("Headers not matching")
	// }

	if actualStatus.Code != expectedStatus.Code {
		t.Errorf("Got Code %s when Code %s was expected", strconv.Itoa(actualStatus.Code), strconv.Itoa(expectedStatus.Code))
	}
}

func TestWriteResponse(t *testing.T) {
	data := []byte("{\"data\":\"Hi\"}")
	status := Status{Code: 204}
	headers := Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Content-Length", "11")

	actualResponse := WriteResponse(data, status, headers)
	expectedResponse := []byte("HTTP/1.1 204 No Content\r\nContent-Type:application/json\r\nContent-Length:11\r\nServer:HomeCloud/0.0.1\r\n\r\n{\"data\":\"Hi\"}")

	if string(actualResponse) != string(expectedResponse) {
		t.Errorf("Got %s, expected %s", string(actualResponse), string(expectedResponse))
	}
}

func TestWritePostRequest(t *testing.T) {
	data := Body{
		Data: "This is a test file.",
	}

	headers := Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Content-Length", "37")

	actualRequest := WritePostRequest("testfile.txt", headers, data)
	expectedRequest := []byte("POST /testfile.txt HTTP/1.1\r\nContent-Type:application/json\r\nContent-Length:37\r\n\r\n{\"data\":\"This is a test file.\"}")

	if string(actualRequest) != string(expectedRequest) {
		t.Errorf("Got %s, expected %s", string(actualRequest), string(expectedRequest))
	}
}

func TestWritePutRequest(t *testing.T) {
	data := Body{
		Data: "This is a put file.",
	}

	headers := Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Content-Length", "36")

	actualRequest := WritePutRequest("testfile.txt", headers, data)

	expectedRequest := []byte("PUT /testfile.txt HTTP/1.1\r\nContent-Type:application/json\r\nContent-Length:36\r\n\r\n{\"data\":\"This is a put file.\"}")

	if string(actualRequest) != string(expectedRequest) {
		t.Errorf("Got %s, expected %s", string(actualRequest), string(expectedRequest))
	}
}
