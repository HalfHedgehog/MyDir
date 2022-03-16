package enum

type ResponseCode int

const (
	Success   ResponseCode = 2000
	Error     ResponseCode = 5000
	Exception ResponseCode = 5001
)
