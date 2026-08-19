package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "book_loop/docs"
)

func New(ctx context.Context, pool *pgxpool.Pool) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	api := engine.Group("/api/v1")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	api.GET("/health", func(c *gin.Context) {
		if err := pool.Ping(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db unreachable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine, nil
}
