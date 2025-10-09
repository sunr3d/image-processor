package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/image-processor/internal/interfaces/services"
)

type Handler struct {
	svc services.ImageService
}

func New(svc services.ImageService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New("")
	router.Use(ginext.Logger(), ginext.Recovery())

	// API
	router.POST("/upload", h.uploadImage)
	router.GET("/image/:id", h.getImage)
	router.DELETE("/image/:id", h.deleteImage)

	// Web-UI
	/* router.Static("/web", "./web")
	router.GET("/", func(c *ginext.Context) {
		c.File("./web/index.html")
	}) */

	return router
}
