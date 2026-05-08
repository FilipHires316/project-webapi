package project_evidence

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"slices"
)

type implPatientEvidenceAPI struct {
}

func NewPatientEvidenceApi() PatientEvidenceAPI {
	return &implPatientEvidenceAPI{}
}

// GET /evidence/{ambulanceId}/patients
func (o implPatientEvidenceAPI) GetEvidedPatients(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		result := ambulance.Patients
		if result == nil {
			result = []EvidedPatient{}
		}
		return nil, result, http.StatusOK
	})
}

// POST /evidence/{ambulanceId}/patients
func (o implPatientEvidenceAPI) CreateEvidedPatient(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		var patient EvidedPatient

		if err := c.ShouldBindJSON(&patient); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if patient.Id == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Patient id is required",
			}, http.StatusBadRequest
		}

		// Skontroluj duplicitu - pacient s rovnakým id už existuje
		conflictIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patient.Id || p.RodneCislo == patient.RodneCislo
		})

		if conflictIndex >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Patient already exists",
			}, http.StatusConflict
		}

		// Inicializuj prázdne pole predpisov ak nie je nastavené
		if patient.Prescriptions == nil {
			patient.Prescriptions = []Prescription{}
		}

		ambulance.Patients = append(ambulance.Patients, patient)

		// Nájdi pacienta späť, aby sme vrátili to čo je teraz v ambulancii
		patientIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patient.Id
		})

		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save patient",
			}, http.StatusInternalServerError
		}

		return ambulance, ambulance.Patients[patientIndex], http.StatusOK
	})
}

// GET /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) GetEvidedPatient(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")

		patientIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patientId
		})

		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		return nil, ambulance.Patients[patientIndex], http.StatusOK
	})
}

// PUT /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) UpdateEvidedPatient(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		var patient EvidedPatient

		if err := c.ShouldBindJSON(&patient); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		patientId := c.Param("patientId")

		patientIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patientId
		})

		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		// Zachovaj id z URL (klient ho v body nemusí mať alebo môže byť nesprávne)
		patient.Id = patientId

		// Zachovaj existujúce predpisy ak v requeste neboli (PUT na pacientovi by ich nemal mazať)
		if patient.Prescriptions == nil {
			patient.Prescriptions = ambulance.Patients[patientIndex].Prescriptions
		}

		ambulance.Patients[patientIndex] = patient

		return ambulance, ambulance.Patients[patientIndex], http.StatusOK
	})
}

// DELETE /evidence/{ambulanceId}/patients/{patientId}
func (o implPatientEvidenceAPI) DeleteEvidedPatient(c *gin.Context) {
	updateAmbulanceFunc(c, func(
		c *gin.Context,
		ambulance *Ambulance,
	) (updatedAmbulance *Ambulance, responseContent interface{}, status int) {
		patientId := c.Param("patientId")

		patientIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patientId
		})

		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		ambulance.Patients = append(ambulance.Patients[:patientIndex], ambulance.Patients[patientIndex+1:]...)

		return ambulance, nil, http.StatusNoContent
	})
}