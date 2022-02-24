package controller

import (
	"myapp/internal/domain/data/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	usecase service.Service
}

func New(service service.Service) Controller {
	return Controller{
		usecase: service,
	}
}

func (c Controller) Init() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "./static")

	router.GET("/", c.HomePage)
	router.GET("/api", c.GetRqstApi)

	return router
}
func (c *Controller) HomePage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "status_page.html", nil)
}
