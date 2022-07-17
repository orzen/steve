package plugin

import "fmt"

type InitHook func() error
type PreReqHook func() error
type PostReqHook func() error

type Plugin interface {
	Init() error
	PreRequest() error
	PostRequest() error
}

func main() {
	fmt.Println("vim-go")
}
