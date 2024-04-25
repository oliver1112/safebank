package domain

type Response struct {
	Status   int64
	Data     interface{}
	ErrorMsg string
}
