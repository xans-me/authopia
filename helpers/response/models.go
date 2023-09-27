package response

import "google.golang.org/genproto/googleapis/rpc/code"

// Struct model response
type Struct struct {
	Data   interface{} `json:"result"`
	TimeIn string      `json:"timeIn"`
}

// ErrorStruct model response
type ErrorStruct struct {
	ErrorCode   code.Code `json:"errorCode"`
	Description string    `json:"description"`
	Message     string    `json:"message"`
}

// we implement the built-in package 'error' interface by creating this function
func (e ErrorStruct) Error() string {
	return e.Message
}
