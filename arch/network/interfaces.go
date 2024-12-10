package network

type ApiError interface {
	GetCode() int
	GetMessage() string
	Error() string
	Unwrap() error
}

type Response interface {
	GetStatus() int
	GetMessage() string
	GetData() any
}
