package httphelper

import (
	"fmt"
	"strings"
)

type EndPoint struct {
	GetEndpoints    map[string]func(Request) []byte
	PostEndpoints   map[string]func(Request) []byte
	PutEndpoints    map[string]func(Request) []byte
	DeleteEndPoints map[string]func(Request) []byte
}

func (e *EndPoint) Get(endpoint string, fn func(Request) []byte) {
	if e.GetEndpoints == nil {
		e.GetEndpoints = make(map[string]func(Request) []byte)
	}

	e.GetEndpoints[endpoint] = fn
}

func (e *EndPoint) Post(endpoint string, fn func(Request) []byte) {
	if e.PostEndpoints == nil {
		e.PostEndpoints = make(map[string]func(Request) []byte)
	}

	e.PostEndpoints[endpoint] = fn
}

func (e *EndPoint) Put(endpoint string, fn func(Request) []byte) {
	if e.PutEndpoints == nil {
		e.PutEndpoints = make(map[string]func(Request) []byte)
	}

	e.PutEndpoints[endpoint] = fn
}

func (e *EndPoint) Delete(endpoint string, fn func(Request) []byte) {
	if e.DeleteEndPoints == nil {
		e.DeleteEndPoints = make(map[string]func(Request) []byte)
	}

	e.DeleteEndPoints[endpoint] = fn
}

func (e EndPoint) Action(method string, uri string) (func(Request) []byte, string) {
	var fn func(Request) []byte
	var key string
	switch strings.ToLower(method) {
	case "get":
		if e.GetEndpoints != nil {
			key = e.ClosestEndpoint(method, uri)
			fmt.Println("Matched GET endpoint: ", key)
			fn = e.GetEndpoints[key]
		}
	case "post":
		if e.PostEndpoints != nil {
			key = e.ClosestEndpoint(method, uri)
			fn = e.PostEndpoints[key]
		}
	case "put":
		if e.PutEndpoints != nil {
			key = e.ClosestEndpoint(method, uri)
			fn = e.PutEndpoints[key]
		}
	case "delete":
		if e.DeleteEndPoints != nil {
			key = e.ClosestEndpoint(method, uri)
			fn = e.DeleteEndPoints[key]
		}
	}
	return fn, key
}

func (e EndPoint) ClosestEndpoint(method string, uri string) string {
	var Closest string
	var Iter map[string]func(Request) []byte
	switch strings.ToLower(method) {
	case "get":
		Iter = e.GetEndpoints
	case "post":
		Iter = e.PostEndpoints
	case "put":
		Iter = e.PutEndpoints
	case "delete":
		Iter = e.DeleteEndPoints
	}

	fmt.Println("Requested resource is ", uri)
	for k := range Iter {
		if uri == k {
			Closest = k
			break
		} else {
			Closest = comparePath(Closest, k, uri)
		}
	}
	fmt.Println("Closest matched endpoiddnt is ", Closest)
	return Closest
}
