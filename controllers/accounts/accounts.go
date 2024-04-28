package accounts

import "github.com/gin-gonic/gin"

type handler struct{}

type Handler interface {
	RouteGroup(*gin.Engine)
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/accounts")

	rg.POST("", h.create)
	rg.GET("/:id", h.get)
}

func (h *handler) create(c *gin.Context) {
	return
}

func (h *handler) get(c *gin.Context) {
	return
}
