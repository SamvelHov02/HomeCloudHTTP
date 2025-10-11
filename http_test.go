package httphelper

import (
	"fmt"
	"strconv"
	"testing"
)

func TestReadResponse(t *testing.T) {
	response := []byte("HTTP/1.1 200 OK\nContent-Type: application/json\nContent-Length: 11\n\n{\"data\":\"Hi\"}")
	actualResponse, _, actualStatus := ReadResponse(response)
	expectedResponse := ResponseBody{
		Data: "Hi",
	}
	fmt.Println("Does it even get here")
	//expectedHeaders := Header{
	//"Content-Type":   []string{"application/json"},
	//"Content-Length": []string{"11"},
	//}

	expectedStatus := HTTPStatus{Code: 200}

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
