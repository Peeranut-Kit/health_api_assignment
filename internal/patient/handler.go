package patient

import (
	"errors"
	"net/http"
	"strconv"

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

// SearchPatient godoc
// @Summary Search for a patient
// @Description Search for a patient which belongs to the same hospital as the staff member in the system
// @Tags Patient
// @Accept json
// @Produce json
// @Param request body pkg.Patient true "Patient search criteria"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /patient/search [get]
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
		return
	}

	// Set search criteria to be the same hospital as current staff
	patientSearchRequest.HospitalID = hospitalIDInt

	// Call service
	patientList, err := h.service.SearchPatient(&patientSearchRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if len(patientList) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No patient found.",
			"data":    patientList,
		})
		return
	}

	// Success searching
	c.JSON(http.StatusOK, gin.H{
		"message": "Search successfully.",
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
	hospitalIDStr, ok := hospitalID.(string)
	if !ok {
		return -1, errors.New("assertion failed for string")
	}

	hospitalIDInt, err := strconv.Atoi(hospitalIDStr)
	if err != nil {
		return -1, errors.New("error converting string to int")
	}

	// Got hospital ID
	return hospitalIDInt, nil
}
