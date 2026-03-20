package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Marosko123/napmap-webapi/api"
	"github.com/Marosko123/napmap-webapi/internal/db_service"
	"github.com/Marosko123/napmap-webapi/internal/napmap"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("NAPMAP_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("NAPMAP_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") {
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update middleware
	dbService := db_service.NewMongoService[napmap.Station](db_service.MongoServiceConfig{})
	defer dbService.Disconnect(context.Background())
	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	// request routings
	handleFunctions := &napmap.ApiHandleFunctions{
		StationsAPI: napmap.NewStationsApi(),
	}
	napmap.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
