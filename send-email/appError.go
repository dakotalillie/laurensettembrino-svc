package main

type codedError struct {
	Code    int
	Message string
}

func (e *codedError) Error() string {
	return e.Message
}
