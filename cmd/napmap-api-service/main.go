package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	// liveness probe — process is up
	engine.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// readiness probe — process can serve traffic (db reachable)
	engine.GET("/ready", func(ctx *gin.Context) {
		probeCtx, cancel := context.WithTimeout(ctx.Request.Context(), 3*time.Second)
		defer cancel()
		if err := dbService.Ping(probeCtx); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// request routings
	handleFunctions := &napmap.ApiHandleFunctions{
		StationsAPI: napmap.NewStationsApi(),
	}
	napmap.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)

	// graceful shutdown — wait for in-flight requests on SIGTERM
	srv := &http.Server{Addr: ":" + port, Handler: engine}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, draining connections")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
