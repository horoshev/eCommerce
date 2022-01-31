package core

import "encoding/json"

type Response struct {
	IsSuccess bool
	Message   string
	Payload   interface{}
}

func Success(message string) Response {
	return Response{IsSuccess: true, Message: message}
}

func Error(message string) Response {
	return Response{IsSuccess: false, Message: message}
}

func NewSuccess(message string, payload interface{}) *Response {
	return &Response{IsSuccess: true, Message: message, Payload: payload}
}

func NewFailure(message string) *Response {
	return &Response{IsSuccess: false, Message: message}
}

func NewError(err error, payload interface{}) *Response {
	return &Response{IsSuccess: false, Message: err.Error(), Payload: payload}
}

func (r Response) Marshal() ([]byte, error) {
	return json.Marshal(r.Payload)
}

func (r Response) StatusHeader() []byte {
	if r.IsSuccess {
		return []byte{1}
	}

	return []byte{0}
}
