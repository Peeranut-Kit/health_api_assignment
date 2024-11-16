package staff_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Peeranut-Kit/health_api_assignment/internal/staff"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock StaffService
type MockStaffService struct {
	mock.Mock
}

func (m *MockStaffService) CreateStaff(staff *pkg.Staff) (*pkg.Staff, error) {
	args := m.Called(staff)
	if args.Get(0) == nil {
		// If the first return value is nil, avoid type assertion and return nil
		return nil, args.Error(1)
	}
	// Type assertion if not nil
	return args.Get(0).(*pkg.Staff), args.Error(1)
}

func (m *MockStaffService) SignInStaff(staff *pkg.Staff) (string, error) {
	args := m.Called(staff)
	return args.String(0), args.Error(1)
}

// Test the CreateStaff handler of HttpStaffrHandler
func TestStaffHandler_CreateStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := staff.NewHttpStaffHandler(mockService)

	r := gin.Default()
	r.POST("/staff/create", handler.CreateStaff)

	// Test case: Successful staff creation
	t.Run("successful staff creation", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		mockService.On("CreateStaff", mock.AnythingOfType("*pkg.Staff")).Return(&inputStaff, nil)

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Equal(t, "Created successfully", response["message"])
		
		mockService.AssertCalled(t, "CreateStaff", &inputStaff)
		// Verify expectations
		mockService.AssertExpectations(t)
	})

	// Test case: Invalid request body format (wrong struct format)
	t.Run("failed staff creation (invalid request body format)", func(t *testing.T) {
		// mock input body request
		inputStaff := "This is random string that should trigger EOF error"

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Failed - Invalidated request body (username or password is empty)
	t.Run("failed staff creation (username or password is empty)", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username:   "",
			Password:   "",
			HospitalID: 1,
		}

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Staff service returns error
	t.Run("staff service error", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		// Reset expectations for this test case
		mockService.ExpectedCalls = nil
		mockService.On("CreateStaff", mock.AnythingOfType("*pkg.Staff")).Return(nil, errors.New("service error"))
		
		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/create", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "service error", response["error"])
		
		mockService.AssertCalled(t, "CreateStaff", &inputStaff)
		// Verify expectations
		mockService.AssertExpectations(t)
	})
}

// Test the SignInStaff handler of HttpStaffrHandler
func TestStaffHandler_SignInStaff(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockStaffService)
	handler := staff.NewHttpStaffHandler(mockService)

	r := gin.Default()
	r.POST("/staff/login", handler.SignInStaff)

	// Test case: Successful staff login
	t.Run("successful staff login", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "secure_password",
		}

		mockService.On("SignInStaff", mock.AnythingOfType("*pkg.Staff")).Return("token", nil)

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "Login successful", response["message"])
		
		mockService.AssertCalled(t, "SignInStaff", &inputStaff)
		// Verify expectations
		mockService.AssertExpectations(t)
	})

	// Test case: Failed - Invalid request body format (wrong struct format)
	t.Run("failed staff login (invalid request body format)", func(t *testing.T) {
		// mock input body request
		inputStaff := "This is random string that should trigger EOF error"

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Failed - Invalidated request body (username or password is empty)
	t.Run("failed staff login (username or password is empty)", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "",
			Password: "",
		}

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Failed - Staff service returns an unauthorized error
	t.Run("unauthorized login attempt", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "wrong_password",
		}

		// Reset expectations for this test case
		mockService.ExpectedCalls = nil
		mockService.On("SignInStaff", mock.AnythingOfType("*pkg.Staff")).Return("", staff.ErrUnauthorized)

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		mockService.AssertCalled(t, "SignInStaff", &inputStaff)
		// Verify expectations
		mockService.AssertExpectations(t)
	})

	// Test case: Staff service Failed - service returns error
	t.Run("staff service error", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "secure_password",
		}

		// Reset expectations for this test case
		mockService.ExpectedCalls = nil
		mockService.On("SignInStaff", mock.AnythingOfType("*pkg.Staff")).Return("", errors.New("service error"))

		body, _ := json.Marshal(inputStaff)
		req := httptest.NewRequest("POST", "/staff/login", bytes.NewBufferString(string(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "service error", response["error"])
		
		mockService.AssertCalled(t, "SignInStaff", &inputStaff)
		// Verify expectations
		mockService.AssertExpectations(t)
	})
}
