package main

import (
	"log"

	"time"

	"github.com/donajivt/go-auth-service/config"
	"github.com/donajivt/go-auth-service/controllers"
	"github.com/donajivt/go-auth-service/db"
	"github.com/donajivt/go-auth-service/middleware"
	"github.com/donajivt/go-auth-service/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	if err := db.Init(); err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	r := gin.Default()

	// Middleware de CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	jwtSvc := services.NewJwtService()
	authSvc := services.NewAuthService(jwtSvc)
	authCtl := controllers.NewAuthController(authSvc)

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authCtl.Register)
			auth.POST("/login", authCtl.Login)
			auth.POST("/assignrole", authCtl.AssignRole)
		}
		// Ruta protegida de prueba
		v1.GET("/ping", middleware.JWTAuth(), func(c *gin.Context) {
			c.JSON(200, gin.H{"msg": "pong"})
		})
	}

	r.Run() // Default port :8080
}
