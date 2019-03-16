package tools

import (
	"context"
	"math/rand"
	"strconv"
	"time"
)

const CTX_LOG_ID = "log_id"

func MakeLogID() string {
	return time.Now().Format("20060102150405") + strconv.Itoa(int(GetLocalIpInt())) + strconv.Itoa(rand.Intn())
}

func NewCtxWithLogID() context.Context {
	ctx := context.Background()
	context.WithValue(ctx, CTX_LOG_ID, MakeLogID())
	return ctx
}
