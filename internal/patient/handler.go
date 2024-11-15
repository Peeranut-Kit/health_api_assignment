package patient

import (
	"errors"
	"net/http"

	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/gin-gonic/gin"
)

// Primary adapter
type PatientHandler struct {
	service PatientServiceInterface
}

// Just define what struct will do
type PatientHandlerInterface interface {
	SearchPatient(c *gin.Context)
}

func NewHttpPatientHandler(service PatientServiceInterface) *PatientHandler {
	return &PatientHandler{service: service}
}

/* This API should mocking searching system of Hospital Information Systems (HIS) of Hospital A API:
Route: GET https://hospital-a.api.co.th/patient/search/{id} */

func (h *PatientHandler) SearchPatient(c *gin.Context) {
	var patientSearchRequest pkg.Patient
	if err := c.ShouldBindJSON(&patientSearchRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve hospital_id
	hospitalIDInt, err := getHospitalID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	// Set search criteria to be the same hospital as current staff
	patientSearchRequest.HospitalID = hospitalIDInt

	// Call service
	patientList, err := h.service.SearchPatient(&patientSearchRequest)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Success searching
	c.JSON(http.StatusOK, gin.H{
		"message": "Patient found successfully.",
		"data":    patientList,
	})
}

func getHospitalID(c *gin.Context) (int, error) {
	// Retrieve hospital_id from gin.Context
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		return -1, errors.New("hospital ID not found")
	}

	// Type assert to int if necessary
	hospitalIDInt, ok := hospitalID.(int)
	if !ok {
		return -1, errors.New("invalid hospital ID format")
	}

	// Got hospital ID
	return hospitalIDInt, nil
}
