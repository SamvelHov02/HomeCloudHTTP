package httphelper

import (
	"encoding/json"
	"fmt"
	"os"
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

// TODO : First need to tet that ReadRequest is working properly
// When testing make sure that the file does not exist before running the test
func TestReadPostMethod(t *testing.T) {
	request := Request{
		Method:   "POST",
		Resource: "/testfile.txt",
		Headers:  Header{"Content-Type": []string{"application/json"}, "Content-Length": []string{"37"}},
		Data: Body{
			Data: "This is a test file.",
		},
	}

	actualData, actualStatus, actualHeaders := ReadPostMethod(request)
	expectedData := Body{
		Data: "This is a test file.",
	}

	fmt.Println("Actual Data:", string(actualData))

	Actuall := &Body{}
	json.Unmarshal(actualData, Actuall)

	expectedStatus := Status{Code: 204}
	expectedHeaders := Header{}
	expectedHeaders.Add("Content-Type", "application/json")
	expectedHeaders.Add("Content-Length", "0")

	if actualStatus.Code != expectedStatus.Code {
		t.Errorf("Got Code %s when Code %s was expected", strconv.Itoa(actualStatus.Code), strconv.Itoa(expectedStatus.Code))
	}

	file, _ := os.ReadFile(filePath + "/testfile.txt")

	if string(file) != expectedData.Data {
		t.Errorf("Got %s, expected %s", Actuall.Data, expectedData.Data)
	}

	for _, key := range expectedHeaders.Keys() {
		if val, ok := actualHeaders.Get(key); !ok || val[0] != expectedHeaders[key][0] {
			t.Errorf("Header %s: got %s, expected %s", key, val, expectedHeaders[key][0])
		}
	}
}
