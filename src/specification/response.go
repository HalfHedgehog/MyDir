package specification

import (
	"httpServer/src/enum"
)

type Response struct {
	Flag    bool              `json:"Flag"`
	Code    enum.ResponseCode `json:"Code"`
	Message string            `json:"Message"`
	Data    interface{}       `json:"Data"`
}

func (*Response) Success(message string, data interface{}) *Response {
	r := Response{
		Flag:    true,
		Code:    enum.Success,
		Message: message,
		Data:    data,
	}
	return &r
}

func (*Response) DefaultSuccess() *Response {
	r := Response{
		Flag:    true,
		Code:    enum.Success,
		Message: "成功！",
		Data:    nil,
	}
	return &r
}

func (*Response) Error(message string, data interface{}) *Response {
	r := Response{
		Flag:    false,
		Code:    enum.Error,
		Message: message,
		Data:    data,
	}
	return &r
}

func (*Response) DefaultError() *Response {
	r := Response{
		Flag:    false,
		Code:    enum.Error,
		Message: "内部错误！",
		Data:    nil,
	}
	return &r
}

func (*Response) Exception(message string, data interface{}) *Response {
	r := Response{
		Flag:    false,
		Code:    enum.Exception,
		Message: message,
		Data:    data,
	}
	return &r
}

func (*Response) CreateResponse(flag bool, code int, message string, data interface{}) *Response {
	r := Response{
		Flag:    flag,
		Code:    enum.ResponseCode(code),
		Message: message,
		Data:    data,
	}
	return &r
}
