package handler

import (
	"github.com/gin-gonic/gin"
	"nsq-demoset/app/app-services/cmd/front_api/middleware"
	"nsq-demoset/app/app-services/ds"
	"nsq-demoset/app/app-services/model"
	"nsq-demoset/app/app-services/repository"
	"nsq-demoset/app/app-services/service"
)

type Handler struct {
	R *gin.Engine

	userScv  model.UserService
	tokenScv model.TokenService
	postScv  model.PostService

	marketService *service.MarketService
	socketService *service.SocketService
}

type HConfig struct {
	R             *gin.Engine
	DS            *ds.DataSource
	MarketRPCAddr string
}

func NewHandler(c *HConfig) *Handler {

	// token
	tokenRepo := repository.NewTokenRepository(c.DS)
	tokenService := service.NewTokenService(&service.TokenConfig{
		TokenRepo: tokenRepo,
	})

	// user repo
	userRepo := repository.NewUserRepository(c.DS)
	userService := service.NewUserService(&service.UserConfig{
		UserRepo: userRepo,
	})

	// post repo
	postRepo := repository.NewPostRepository(c.DS)
	postService := service.NewPostService(&service.PostConfig{
		PostRepo: postRepo,
	})

	// market service
	marketService := service.NewMarketService(c.MarketRPCAddr)

	// socket service
	socketService := service.NewSocketService(&service.SConfig{
		MarketService: marketService,
	})

	return &Handler{
		R:        c.R,
		userScv:  userService,
		tokenScv: tokenService,
		postScv:  postService,

		marketService: marketService,
		socketService: socketService,
	}
}

func (h *Handler) Register() {
	// register cors middleware
	h.R.Use(middleware.Cors())

	// home
	homeHandler := NewHomeHandler(h)
	homeHandler.Register()

	// auth
	authHandler := NewAuthHandler(h)
	authHandler.Register()

	// post
	postHandler := NewPostHandler(h)
	postHandler.Register()

	// web-socket
	wsHandler := NewWSHandler(h)
	wsHandler.Register()

}
