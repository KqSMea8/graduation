package tools

import "context"

const CTX_LOG_ID = "log_id"

func MakeLogID() string {
	
}

func NewCtxWithID() context.Context {
	ctx := context.Background()
	context.WithValue(ctx, CTX_LOG_ID, )
}
