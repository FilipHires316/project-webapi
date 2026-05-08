package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/FilipHires316/project-webapi/api"
    "github.com/FilipHires316/project-webapi/internal/project_evidence"
    "github.com/FilipHires316/project-webapi/internal/db_service"
    "context"
    "time"
    "github.com/gin-contrib/cors"
)

func main() {
    log.Printf("Server started")
    port := os.Getenv("PROJECT_API_PORT")
    if port == "" {
        port = "8080"
    }
    environment := os.Getenv("PROJECT_API_ENVIRONMENT")
    if !strings.EqualFold(environment, "production") { // case insensitive comparison
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
        MaxAge: 12 * time.Hour,
    })
    engine.Use(corsMiddleware)
        // setup context update  middleware
    dbService := db_service.NewMongoService[project_evidence.Ambulance](db_service.MongoServiceConfig{})
    defer dbService.Disconnect(context.Background())
    engine.Use(func(ctx *gin.Context) {
        ctx.Set("db_service", dbService)
        ctx.Next()
    })
    // request routings
    handleFunctions := &project_evidence.ApiHandleFunctions{
		PatientEvidenceAPI:      project_evidence.NewPatientEvidenceApi(),
		PatientPrescriptionsAPI: project_evidence.NewPatientPrescriptionsApi(),
        AmbulancesAPI:           project_evidence.NewAmbulancesApi(),
	}
	project_evidence.NewRouterWithGinEngine(engine, *handleFunctions)
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}