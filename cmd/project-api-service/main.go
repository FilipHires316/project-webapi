package main

import (
    "log"
    "os"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/FilipHires316/project-webapi/api"
    "github.com/FilipHires316/project-webapi/internal/project_evidence"
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
    // request routings
    handleFunctions := &project_evidence.ApiHandleFunctions{
		PatientEvidenceAPI:      project_evidence.NewPatientEvidenceApi(),
		PatientPrescriptionsAPI: project_evidence.NewPatientPrescriptionsApi(),
	}
	project_evidence.NewRouterWithGinEngine(engine, *handleFunctions)
    engine.GET("/openapi", api.HandleOpenApi)
    engine.Run(":" + port)
}