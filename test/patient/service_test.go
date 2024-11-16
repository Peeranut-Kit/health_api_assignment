package patient

import (
	"errors"
	"testing"

	"github.com/Peeranut-Kit/health_api_assignment/internal/patient"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPatientRepo struct {
	mock.Mock
}

func (m *mockPatientRepo) SearchPatient(request *pkg.Patient) ([]pkg.Patient, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		// If the first return value is nil, avoid type assertion and return nil
		return nil, args.Error(1)
	}
	// Type assertion if not nil
	return args.Get(0).([]pkg.Patient), args.Error(1)
}

func TestPatientService_SearchPatient(t *testing.T) {
	mockRepo := new(mockPatientRepo)
	service := patient.NewPatientService(mockRepo)

	// Test case: Successful patient searching
	t.Run("successful patient searching", func(t *testing.T) {
		// mock input body request
		inputPatient := pkg.Patient{
			PatientHN:  "654350968",
			HospitalID: 1,
		}

		mockRepo.On("SearchPatient", &inputPatient).Return([]pkg.Patient{{ID: 1, FirstNameEn: "John"}}, nil)

		paientList, err := service.SearchPatient(&inputPatient)

		assert.NoError(t, err)
		assert.NotEmpty(t, paientList)
		
		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case: Failed - patient searching error
	t.Run("error patient repository searching", func(t *testing.T) {
		// mock input body request
		inputPatient := pkg.Patient{
			PatientHN:  "-7",
			HospitalID: 1,
		}

		mockRepo.On("SearchPatient", &inputPatient).Return(nil, errors.New("database error"))

		_, err := service.SearchPatient(&inputPatient)

		assert.Error(t, err)
		assert.EqualError(t, err, "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}
