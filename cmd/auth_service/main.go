package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/tour/internal/config"
	"gitlab.com/tour/internal/core/repository"
	"gitlab.com/tour/internal/pkg/logger"
	redis_pkg "gitlab.com/tour/internal/pkg/redis"
	"gitlab.com/tour/internal/transport/grpc"
	"gitlab.com/tour/internal/transport/handlers"
	grpcMain "google.golang.org/grpc"
)

var (
	cfg *config.Config
)

func init() {
	log.Println("Initializing...")
	cfg = config.Load()
	log.Println("Initializing done")
}

func main() {

	_logger := logger.SetLogger(&cfg.Logging)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Project.Timeout)
	defer cancel()

	repos := repository.New(ctx, cfg.PSQL.URI)

	redisClient, err := redis_pkg.InitRedis(cfg.RedisAddress)
	if err != nil {
		log.Fatal("Error while connecting to redis: ", err)
		return
	}

	grpcServer := grpc.New(repos, redisClient, cfg, _logger)

	lis, err := net.Listen("tcp", cfg.Grpc.URL)
	if err != nil {
		log.Fatal("Error while listening: ", err)
		return
	}

	go func() {
		logger.Log.Info("starting grpc server on " + cfg.Grpc.URL)

		err := grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	router := setUpHttp()
	go func() {
		logger.Log.Info("starting http server on " + cfg.Http.URL)
		router.Run(cfg.Http.URL)
	}()

	gracefulShutdown(grpcServer, ctx, cancel)
}

func gracefulShutdown(grpcServer *grpcMain.Server, ctx context.Context, cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		logger.Log.Info("shutting down")

		grpcServer.GracefulStop()

		logger.Log.Info("shutdown successfully called")
		wg.Done()
	}(&wg)

	go func() {
		wg.Wait()
		cancel()
	}()

	<-ctx.Done()
	os.Exit(0)
}

func setUpHttp() *gin.Engine {
	router := gin.Default()

	switch cfg.Mode {
	case "dev":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "*")

	router.Use(cors.New(config))
	v1 := router.Group("/v1")

	gwMux := handlers.New(context.Background(), cfg)
	router.Static("/swagger", "./doc/swagger")
	v1.Any("/*any", func(c *gin.Context) {
		gwMux.ServeHTTP(c.Writer, c.Request)
	})

	return router
}
