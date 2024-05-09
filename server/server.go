package server

import (
	"fmt"
	"sync"

	"github.com/DanielVieirass/um_help/config"
	"github.com/DanielVieirass/um_help/server/controller"
	"github.com/DanielVieirass/um_help/server/middleware"
	"github.com/DanielVieirass/um_help/server/router"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var (
	instance *Server
	once     sync.Once
)

type Server struct {
	cfg    *config.Config
	svr    *echoadapter.EchoLambda
	logger *zerolog.Logger
	ctrl   *controller.Controller
}

func New(cfg *config.Config, logger *zerolog.Logger, ctrl *controller.Controller) *Server {
	once.Do(func() {
		svr := echo.New()

		svr.HideBanner = true
		svr.HidePort = true

		middleware.SetMiddlewares(svr, cfg)
		router.Register(cfg, svr, ctrl)

		instance = &Server{
			cfg:    cfg,
			svr:    echoadapter.New(svr),
			logger: logger,
			ctrl:   ctrl,
		}
	})

	return instance
}

func (s *Server) Start() error {
	s.logger.Info().Msg("starting server")

	if err := s.svr.Echo.Start(fmt.Sprintf(":%d", s.cfg.InternalConfig.ServerPort)); err != nil {
		return err
	}

	return nil
}
