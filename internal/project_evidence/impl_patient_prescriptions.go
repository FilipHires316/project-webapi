package project_evidence

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPatientPrescriptionsAPI struct {
}

func NewPatientPrescriptionsApi() PatientPrescriptionsAPI {
	return &implPatientPrescriptionsAPI{}
}

// GET /evidence/{ambulanceId}/patients/{patientId}/prescriptions
func (o implPatientPrescriptionsAPI) GetPatientPrescriptions(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// POST /evidence/{ambulanceId}/patients/{patientId}/prescriptions
func (o implPatientPrescriptionsAPI) CreatePrescription(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// GET /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) GetPrescription(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// PUT /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) UpdatePrescription(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

// DELETE /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) DeletePrescription(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}