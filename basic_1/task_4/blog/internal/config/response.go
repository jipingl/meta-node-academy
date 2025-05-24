package config

const (
	CodeSuccess       = 200
	CodeInternalError = 500
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) Resp {
	return Resp{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

func Fail(code int, err error) Resp {
	return Resp{
		Code:    code,
		Message: err.Error(),
		Data:    nil,
	}
}
