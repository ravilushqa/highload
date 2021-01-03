package tarantool

import (
	"fmt"

	"github.com/tarantool/go-tarantool"
)

func New(addr, user, pass string) (*tarantool.Connection, error) {
	if addr == "" {
		fmt.Println("empty tarantool addr")
		return nil, nil
	}
	opts := tarantool.Opts{
		User: user,
		Pass: pass,
	}
	return tarantool.Connect(addr, opts)
}
