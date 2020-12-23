package tarantool

import "github.com/tarantool/go-tarantool"

func New(addr, user, pass string) (*tarantool.Connection, error) {
	opts := tarantool.Opts{
		User: user,
		Pass: pass,
	}
	return tarantool.Connect(addr, opts)
}
