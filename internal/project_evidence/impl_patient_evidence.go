package project_evidence

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type implPatientEvidenceAPI struct {
}

func NewPatientEvidenceApi() PatientEvidenceAPI {
	return &implPatientEvidenceAPI{}
}

// validateEvidedPatient kontroluje povinné polia a enum hodnoty.
// Vracia gin.H s chybou a status kódom, alebo nil ak je všetko v poriadku.
func validateEvidedPatient(patient EvidedPatient) (gin.H, int) {
	missing := validateRequired(map[string]string{
		"id":          patient.Id,
		"name":        patient.Name,
		"rodneCislo":  patient.RodneCislo,
		"dateOfBirth": patient.DateOfBirth,
		"gender":      patient.Gender,
		"insurance":   patient.Insurance,
	})
	if len(missing) > 0 {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Missing required fields",
			"missing": missing,
		}, http.StatusBadRequest
	}

	validGenders := []string{"male", "female", "other"}
	if !slices.Contains(validGenders, patient.Gender) {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid gender, must be one of: male, female, other",
		}, http.StatusBadRequest
	}

	validInsurances := []string{"VšZP", "Dôvera", "Union"}
	if !slices.Contains(validInsurances, patient.Insurance) {
		return gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid insurance, must be one of: VšZP, Dôvera, Union",
		}, http.StatusBadRequest
	}

	if patient.BloodType != "" {
		validBloodTypes := []string{"A+", "A-", "B+", "B-", "AB+", "AB-", "0+", "0-"}
		if !slices.Contains(validBloodTypes, patient.BloodType) {
			return gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid bloodType, must be one of: A+, A-, B+, B-, AB+, AB-, 0+, 0-",
			}, http.StatusBadRequest
		}
	}

	return nil, 0
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

		// Validácia povinných polí a enum hodnôt
		if errResp, errStatus := validateEvidedPatient(patient); errResp != nil {
			return nil, errResp, errStatus
		}

		// Skontroluj duplicitu - pacient s rovnakým id alebo rodným číslom už existuje
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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}

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
		patientId := c.Param("patientId")

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}

		var patient EvidedPatient
		if err := c.ShouldBindJSON(&patient); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		// Forsuj id z URL pred validáciou - klient ho v body nemusí mať
		// alebo ho môže mať s inou hodnotou. Validátor potom skontroluje že je vyplnené.
		patient.Id = patientId

		// Validácia povinných polí a enum hodnôt (rovnako ako pri create)
		if errResp, errStatus := validateEvidedPatient(patient); errResp != nil {
			return nil, errResp, errStatus
		}

		patientIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id == patientId
		})

		if patientIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Patient not found",
			}, http.StatusNotFound
		}

		// Skontroluj že nový rodneCislo nepatrí inému pacientovi
		duplicateRodneCisloIndex := slices.IndexFunc(ambulance.Patients, func(p EvidedPatient) bool {
			return p.Id != patientId && p.RodneCislo == patient.RodneCislo
		})
		if duplicateRodneCisloIndex >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Another patient with this rodneCislo already exists",
			}, http.StatusConflict
		}

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

		if patientId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "patientId path parameter is required",
			}, http.StatusBadRequest
		}

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