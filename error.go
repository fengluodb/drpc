package drpc

import "errors"

var (
	ErrUnmarshal = errors.New("an error occurred in Unmarshal")
	ErrShutdown  = errors.New("connection is shut down")
)
