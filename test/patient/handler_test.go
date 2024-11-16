package patient_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peeranut-Kit/health_api_assignment/internal/patient"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock PatientService
type MockPatientService struct {
	mock.Mock
}

func (m *MockPatientService) SearchPatient(patientSearchRequest *pkg.Patient) ([]pkg.Patient, error) {
	args := m.Called(patientSearchRequest)
	if args.Get(0) == nil {
		// If the first return value is nil, avoid type assertion and return nil
		return nil, args.Error(1)
	}
	// Type assertion if not nil
	return args.Get(0).([]pkg.Patient), args.Error(1)
}

// Mock returning hospitalID as 1 without JWT cookie
func mockGetHospitalID(c *gin.Context) (int, error) {
	return 1, nil
}

// Tests the SearchPatient handler of HttpPatientHandler
func TestPatientHandler_SearchPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockPatientService)
	handler := &patient.PatientHandler{
		Service:         mockService,
		GetHospitalIDFn: mockGetHospitalID,
	}

	r := gin.Default()
	r.GET("/patient/search", handler.SearchPatient)

	// Test case: Successful patient searching
	t.Run("successful patient searching", func(t *testing.T) {
		// mock input body request
		inputPatientSearchRequest := pkg.Patient{
			PatientHN: "654350968",
		}

		// mock SearchPatient
		mockService.On("SearchPatient", mock.AnythingOfType("*pkg.Patient")).Return([]pkg.Patient{{ID: 1, FirstNameEn: "John"}}, nil)

		body, _ := json.Marshal(inputPatientSearchRequest)
		req := httptest.NewRequest("GET", "/patient/search", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Search successfully.", response["message"])

		// Verify expectations
		mockService.AssertExpectations(t)
	})

	// Test case: Successful patient searching but not found
	t.Run("patient searching not found", func(t *testing.T) {
		// mock input body request
		inputPatientSearchRequest := pkg.Patient{
			PatientHN: "82",
		}

		// Reset expectations for this test case
		mockService.ExpectedCalls = nil
		mockService.On("SearchPatient", mock.AnythingOfType("*pkg.Patient")).Return([]pkg.Patient{}, nil)

		body, _ := json.Marshal(inputPatientSearchRequest)
		req := httptest.NewRequest("GET", "/patient/search", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "No patient found.", response["message"])

		// Verify expectations
		mockService.AssertExpectations(t)
	})

	// Test case: Failed - Invalid request body format (wrong struct format)
	t.Run("failed searching login (invalid request body format)", func(t *testing.T) {
		// mock input body request
		inputPatientSearchRequest := "This is random string that should trigger EOF error"

		body, _ := json.Marshal(inputPatientSearchRequest)
		req := httptest.NewRequest("GET", "/patient/search", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Patient service Failed - service returns an error
	t.Run("patient service error", func(t *testing.T) {
		// mock input body request
		inputPatientSearchRequest := pkg.Patient{
			PatientHN: "654350968",
		}

		// Reset expectations for this test case
		mockService.ExpectedCalls = nil
		mockService.On("SearchPatient", mock.AnythingOfType("*pkg.Patient")).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(inputPatientSearchRequest)
		req := httptest.NewRequest("GET", "/patient/search", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "service error", response["error"])

		// Verify expectations
		mockService.AssertExpectations(t)
	})
}
