package client

import "fmt"

var (
	InvalidParam = fmt.Errorf("invalid param")
	DiscardError = fmt.Errorf("order expired, discard it")
)
