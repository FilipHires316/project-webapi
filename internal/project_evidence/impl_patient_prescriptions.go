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

// validatePrescription kontroluje povinné polia a enum hodnoty predpisu.
// Vracia gin.H s chybou a status kódom, alebo nil ak je všetko v poriadku.
func validatePrescription(prescription Prescription) (gin.H, int) {
	missing := validateRequired(map[string]string{
		"id":             prescription.Id,
		"medicineName":   prescription.MedicineName,
		"strength":       prescription.Strength,
		"form":           prescription.Form,
		"dosage":         prescription.Dosage,
		"quantity":       prescription.Quantity,
		"prescribedDate": prescription.PrescribedDate,
		"validUntil":     prescription.ValidUntil,
		"prescribedBy":   prescription.PrescribedBy,
		"status":         prescription.Status,
	})
	if len(missing) > 0 {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Missing required fields",
			"missing": missing,
		}, http.StatusBadRequest
	}

	validForms := []string{"tbl.", "kapsule", "sirup", "kvapky", "masť", "krém", "injekcia", "inhalátor", "čapík"}
	if !slices.Contains(validForms, prescription.Form) {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid form, must be one of: tbl., kapsule, sirup, kvapky, masť, krém, injekcia, inhalátor, čapík",
		}, http.StatusBadRequest
	}

	validStatuses := []string{"active", "dispensed", "expired"}
	if !slices.Contains(validStatuses, prescription.Status) {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid status, must be one of: active, dispensed, expired",
		}, http.StatusBadRequest
	}

	if prescription.Coverage != "" {
		validCoverages := []string{"full", "partial", "none"}
		if !slices.Contains(validCoverages, prescription.Coverage) {
			return gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid coverage, must be one of: full, partial, none",
			}, http.StatusBadRequest
		}
	}

	if prescription.RepeatMonths != 0 {
		validRepeats := []int32{1, 3, 6, 12}
		if !slices.Contains(validRepeats, prescription.RepeatMonths) {
			return gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid repeatMonths, must be one of: 1, 3, 6, 12",
			}, http.StatusBadRequest
		}
	}

	return nil, 0
}

// GET /evidence/{ambulanceId}/patients/{patientId}/prescriptions
func (o implPatientPrescriptionsAPI) GetPatientPrescriptions(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}

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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}

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

		// Validácia povinných polí a enum hodnôt
		if errResp, errStatus := validatePrescription(prescription); errResp != nil {
			return nil, errResp, errStatus
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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}
		if prescriptionId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "prescriptionId path parameter is required",
			}, http.StatusBadRequest
		}

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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}
		if prescriptionId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "prescriptionId path parameter is required",
			}, http.StatusBadRequest
		}

		var prescription Prescription
		if err := c.ShouldBindJSON(&prescription); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		// Forsuj id z URL pred validáciou - klient ho v body nemusí mať
		// alebo ho môže mať s inou hodnotou. Validátor potom skontroluje že je vyplnené.
		prescription.Id = prescriptionId

		// Validácia povinných polí a enum hodnôt (rovnako ako pri create)
		if errResp, errStatus := validatePrescription(prescription); errResp != nil {
			return nil, errResp, errStatus
		}

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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}
		if prescriptionId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "prescriptionId path parameter is required",
			}, http.StatusBadRequest
		}

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