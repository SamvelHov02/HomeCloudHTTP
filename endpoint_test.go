package httphelper

import (
	"testing"
)

func TestClosestEndpoint(t *testing.T) {
	fn1 := func(r Request) []byte {
		return []byte("test")
	}

	fn2 := func(r Request) []byte {
		return []byte("test2")
	}
	endpoint := EndPoint{}
	endpoint.Get("/users/profile", fn1)
	endpoint.Get("/users/profile/samo", fn2)

	endpointKey := endpoint.ClosestEndpoint("GET", "/users/profile/samo/settings")
	if endpointKey != "/users/profile/samo" {
		t.Errorf("Expected /users/profile/samo but got %s", endpointKey)
	}
}

func TestComparePath(t *testing.T) {
	path1 := "/users/profile"
	path2 := "/users/profile/samo"
	reference := "/users/profile/samo/settings"

	expected := "/users/profile/samo"

	result := comparePath(path1, path2, reference)

	if result != expected {
		t.Errorf("Expected %s but got %s", expected, result)
	}
}
