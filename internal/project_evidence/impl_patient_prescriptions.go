package project_evidence

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type implPatientPrescriptionsAPI struct {
}

func NewPatientPrescriptionsApi() PatientPrescriptionsAPI {
	return &implPatientPrescriptionsAPI{}
}

// findPatientIndex je pomocná funkcia ktorá nájde pacienta podľa id v poli pacientov ambulancie.
// Vracia index alebo -1 ak pacient neexistuje.
func findPatientIndex(ambulance *Ambulance, patientId string) int {
	return slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
		return p.Id == patientId
	})
}

// GET /evidence/{ambulanceId}/patients/{patientId}/prescriptions
func (o implPatientPrescriptionsAPI) GetPatientPrescriptions(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")

		patientIndex := findPatientIndex(ambulance, patientId)
		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		result := ambulance.Patients[patientIndex].Prescriptions
		if result == nil {
			result = []Prescription{}
		}

		return nil, result, http.StatusOK
	})
}

// POST /evidence/{ambulanceId}/patients/{patientId}/prescriptions
func (o implPatientPrescriptionsAPI) CreatePrescription(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")

		patientIndex := findPatientIndex(ambulance, patientId)
		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		var prescription Prescription
		if err := c.ShouldBindJSON(&prescription); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if prescription.Id == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Prescription id is required",
			}, http.StatusBadRequest
		}

		// Skontroluj duplicitu - predpis s rovnakým id už existuje
		conflictIndex := slices.IndexFunc(ambulance.Patients[patientIndex].Prescriptions, func(rx Prescription) bool {
			return rx.Id == prescription.Id
		})

		if conflictIndex >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Prescription already exists",
			}, http.StatusConflict
		}

		ambulance.Patients[patientIndex].Prescriptions = append(
			ambulance.Patients[patientIndex].Prescriptions,
			prescription,
		)

		// Nájdi predpis späť aby sme vrátili to čo je teraz uložené
		rxIndex := slices.IndexFunc(ambulance.Patients[patientIndex].Prescriptions, func(rx Prescription) bool {
			return rx.Id == prescription.Id
		})

		if rxIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save prescription",
			}, http.StatusInternalServerError
		}

		return ambulance, ambulance.Patients[patientIndex].Prescriptions[rxIndex], http.StatusOK
	})
}

// GET /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) GetPrescription(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")
		prescriptionId := c.Param("prescriptionId")

		patientIndex := findPatientIndex(ambulance, patientId)
		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		rxIndex := slices.IndexFunc(ambulance.Patients[patientIndex].Prescriptions, func(rx Prescription) bool {
			return rx.Id == prescriptionId
		})

		if rxIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Prescription not found",
			}, http.StatusNotFound
		}

		return nil, ambulance.Patients[patientIndex].Prescriptions[rxIndex], http.StatusOK
	})
}

// PUT /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) UpdatePrescription(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")
		prescriptionId := c.Param("prescriptionId")

		patientIndex := findPatientIndex(ambulance, patientId)
		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		rxIndex := slices.IndexFunc(ambulance.Patients[patientIndex].Prescriptions, func(rx Prescription) bool {
			return rx.Id == prescriptionId
		})

		if rxIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Prescription not found",
			}, http.StatusNotFound
		}

		var prescription Prescription
		if err := c.ShouldBindJSON(&prescription); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		// Zachovaj id z URL aby klient nemohol predpis premenovať
		prescription.Id = prescriptionId

		ambulance.Patients[patientIndex].Prescriptions[rxIndex] = prescription

		return ambulance, ambulance.Patients[patientIndex].Prescriptions[rxIndex], http.StatusOK
	})
}

// DELETE /evidence/{ambulanceId}/patients/{patientId}/prescriptions/{prescriptionId}
func (o implPatientPrescriptionsAPI) DeletePrescription(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")
		prescriptionId := c.Param("prescriptionId")

		patientIndex := findPatientIndex(ambulance, patientId)
		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		rxIndex := slices.IndexFunc(ambulance.Patients[patientIndex].Prescriptions, func(rx Prescription) bool {
			return rx.Id == prescriptionId
		})

		if rxIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Prescription not found",
			}, http.StatusNotFound
		}

		ambulance.Patients[patientIndex].Prescriptions = append(
			ambulance.Patients[patientIndex].Prescriptions[:rxIndex],
			ambulance.Patients[patientIndex].Prescriptions[rxIndex+1:]...,
		)

		return ambulance, nil, http.StatusNoContent
	})
}