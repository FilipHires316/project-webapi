package project_evidence

import (
	"net/http"

	"github.com/FilipHires316/project-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
)

type implAmbulancesAPI struct {
}

func NewAmbulancesApi() AmbulancesAPI {
	return &implAmbulancesAPI{}
}

// CreateAmbulance handles POST /ambulance
func (o implAmbulancesAPI) CreateAmbulance(c *gin.Context) {
	// 1. Získaj DB service z context-u (nastavený v main.go middleware)
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of required type",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	// 2. Rozparsuj telo requestu na Ambulance objekt
	ambulance := Ambulance{}
	err := c.BindJSON(&ambulance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// 3. Validácia povinných polí
	if ambulance.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Ambulance id is required",
		})
		return
	}

	// 4. Inicializuj prázdny zoznam pacientov ak nebol poslaný
	if ambulance.Patients == nil {
		ambulance.Patients = []EvidedPatient{}
	}

	// 5. Vytvor dokument v MongoDB
	err = db.CreateDocument(c, ambulance.Id, &ambulance)

	switch err {
	case nil:
		c.JSON(http.StatusOK, ambulance)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{
			"status":  "Conflict",
			"message": "Ambulance already exists",
			"error":   err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to create ambulance in database",
			"error":   err.Error(),
		})
	}
}

// DeleteAmbulance handles DELETE /ambulance/:ambulanceId
func (o implAmbulancesAPI) DeleteAmbulance(c *gin.Context) {
	// 1. Získaj DB service z context-u
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of required type",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	// 2. Získaj ID ambulancie z URL parametra
	ambulanceId := c.Param("ambulanceId")

	// 3. Zmaž dokument z MongoDB
	err := db.DeleteDocument(c, ambulanceId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Ambulance not found",
			"error":   err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to delete ambulance from database",
			"error":   err.Error(),
		})
	}
}