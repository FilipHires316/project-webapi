package project_evidence

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// implPatientEvidenceAPI implementuje rozhranie PatientEvidenceAPI.
// Štruktúra je zámerne neexportovaná (malé písmeno) — inštanciu vytvárame
// cez NewPatientEvidenceApi(), čo zabezpečí kompilačnú kontrolu, že
// štruktúra rozhranie naozaj implementuje.
type implPatientEvidenceAPI struct {
}

func NewPatientEvidenceApi() PatientEvidenceAPI {
	return &implPatientEvidenceAPI{}
}

// GET /evidence/{ambulanceId}/patients
func (o implPatientEvidenceAPI) GetEvidedPatients(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// POST /evidence/{ambulanceId}/patients
func (o implPatientEvidenceAPI) CreateEvidedPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// GET /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) GetEvidedPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// PUT /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) UpdateEvidedPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// DELETE /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) DeleteEvidedPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}