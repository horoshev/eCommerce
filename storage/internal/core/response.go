package core

import "encoding/json"

type Response struct {
	IsSuccess bool
	Message   string
	Payload   interface{}
}

func NewSuccess(message string, payload interface{}) *Response {
	return &Response{IsSuccess: true, Message: message, Payload: payload}
}

func NewError(err error, payload interface{}) *Response {
	return &Response{IsSuccess: false, Message: err.Error(), Payload: payload}
}

func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r.Payload)
}

func (r *Response) StatusHeader() []byte {
	if r.IsSuccess {
		return []byte{1}
	}

	return []byte{0}
}
