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

	// InternalError error type
	InternalError = 3
)

// Request represents a request
type Request struct {
	Action Verb              `json:"action"`
	Args   map[string]string `json:"args,omitempty"`
}

// Response represents a request
type Response struct {
	Status int    `json:"status"`
	Data 	[]byte 	`json:"data"`
}

// FetchResponse represents the reponse emitted
// when a fetch request has been submitted
type FetchResponse struct {
	Status int `json:"status"`
	Store Store `json:"store"`
}

// Error represents an error
type Error struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
}

// handleRequest
func handleRequest(message []byte) []byte {
	var request Request
	var payload []byte
	var err error

	err = json.Unmarshal(message, &request)
	if err != nil {
		payload, _ = json.Marshal(NewError(InvalidRequest, "unable to deserialize provided json data"))
		return payload
	}

	switch request.Action {
	case Fetch:
		response, err := fetch()
		if err != nil {
			errBytes, _ := json.Marshal(NewError(InternalError, err.Error()))
			return errBytes
		}

		responseBytes, _ := json.Marshal(response)
		return responseBytes
	case Push:
		break
	}

	return payload
}

func fetch() (*FetchResponse, error) {
	storePath, err := findBivouacFile()
	if err != nil {
			return nil, err
	}

	store, err := LoadStore(storePath)
	if err != nil {
		return nil, err
	}

	return NewFetchResponse(Success, *store), nil
}

// NewRequest instanciates a new Request with provided action and args
func NewRequest(action Verb, args map[string]string) *Request {
	return &Request{
		Action: action,
		Args:   args,
	}
}

// NewResponse instanciates a new Response with provided status and data
func NewResponse(status int, data []byte) *Response {
	return &Response{
		Status: status,
		Data:   data,
	}
}

// NewFetchResponse instanciates a new FetchResponse
func NewFetchResponse(status int, store Store) *FetchResponse {
	return &FetchResponse{
		Status: status,
		Store: store,
	}
}

// NewError instanciates a new Error with provided type and message
func NewError(errType int, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}
