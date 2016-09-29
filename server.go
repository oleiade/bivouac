package main


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
)

// Request represents a request
type Request struct {
  Action Verb `json:"action"`
  Args map[string]string `json:"args"`
}

// Response represents a request
type Response struct {
  Status int `json:"status"`
  Data map[string]string `json:"data"`
}

// Error represents an error
type Error struct {
  Type int `json:"type"`
  Message string `json:"message"`
}
