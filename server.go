package main

import "encoding/json"

// Verb represents a request action's verb
type Verb string

const (
	// Fetch verb. used for requesting the store
	// from the server
	Fetch Verb = "fetch"

	// Push verb. Used for pushing a new store
	// content to the server.
	Push = "push"
)

const (
	// Success status code
	Success int = 200

	// Failure status code
	Failure int = 400
)

const (
	// NotFound error type
	NotFound int = 1

	// InvalidRequest error type
	InvalidRequest = 2
)

// Request represents a request
type Request struct {
	Action Verb              `json:"action"`
	Args   map[string]string `json:"args"`
}

// Response represents a request
type Response struct {
	Status int               `json:"status"`
	Data   map[string]string `json:"data"`
}

// Error represents an error
type Error struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

// handleRequest
func handleRequest(message []byte) []byte {
	var request Request
	var err error

	err = json.Unmarshal(message, &request)
	if err != nil {
		message, _ := json.Marshal(NewError(InvalidRequest, "unable to deserialize provided json data"))
		return message
	}

	switch request.Action {
	case Fetch:
		response := NewResponse(Success, make(map[string]string))
		break
	case Push:
		break
	}

	return NewResponse(Success, make(map[string]string))
}

func fetch() *Response {

	return NewResponse(Success)
}

// NewRequest instanciates a new Request with provided action and args
func NewRequest(action Verb, args map[string]string) *Request {
	return &Request{
		Action: action,
		Args:   args,
	}
}

// NewResponse instanciates a new Response with provided status and data
func NewResponse(status int, data map[string]string) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}

// NewError instanciates a new Error with provided type and message
func NewError(errType int, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}
