package handler

import (
	"github.com/g10guang/graduation/store"
)

var storage store.Storage

func init() {
	storage = store.NewLocalStorage()
}
