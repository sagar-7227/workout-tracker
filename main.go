package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"          
    ginSwagger "github.com/swaggo/gin-swagger"                      
	"workout-tracker/config"
	"workout-tracker/routes"
	"workout-tracker/seed"
	_ "workout-tracker/docs"
)

// @title Workout Tracker API
// @version 1.0
// @description REST API for managing workouts, exercises, and reports.
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use format: Bearer {token}
func main() {
	_ = godotenv.Load()

	db := config.MustInitDB()
	seed.RunExerciseSeeder(db)

	r := gin.Default()
	routes.Register(r, db)

	if os.Getenv("APP_ENV") != "prod" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	srv := &http.Server{Addr: ":" + config.GetPort(), Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
