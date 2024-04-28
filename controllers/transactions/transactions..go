package transactions

import "github.com/gin-gonic/gin"

type handler struct{}

type Handler interface {
	RouteGroup(*gin.Engine)
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) RouteGroup(r *gin.Engine) {
	rg := r.Group("/transactions")

	rg.POST("", h.create)
}

func (h *handler) create(c *gin.Context) {
	return
}
