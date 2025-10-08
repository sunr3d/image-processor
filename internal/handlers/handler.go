package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/image-processor/internal/interfaces/services"
)

type Handler struct {
	svc services.ImageProcessor
}

func New(svc services.ImageProcessor) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger(), ginext.Recovery())

	// TODO: Implement handlers

	return router
}
