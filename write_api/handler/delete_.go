package handler

import (
	"context"
	"net/http"
)

type DeleteHandler struct {
	*CommonHandler
}

func NewDeleteHandler(r *http.Request) *DeleteHandler {
	return &DeleteHandler{
		CommonHandler: NewCommonHandler(r),
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	return nil
}

