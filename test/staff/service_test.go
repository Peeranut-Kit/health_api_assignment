package staff_test

import (
	"errors"
	"testing"

	"github.com/Peeranut-Kit/health_api_assignment/internal/staff"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type mockStaffRepo struct {
	mock.Mock
}

func (m *mockStaffRepo) CreateStaff(staff *pkg.Staff) error {
	args := m.Called(staff)
	return args.Error(0)
}

func (m *mockStaffRepo) GetStaffFromUsername(username string) (*pkg.Staff, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		// If the first return value is nil, avoid type assertion and return nil
		return nil, args.Error(1)
	}
	// Type assertion if not nil
	return args.Get(0).(*pkg.Staff), args.Error(1)
}

// Mock bcrypt hasher
type MockPasswordHasher struct {
	mock.Mock
}

// Mock bcrypt.CompareHashAndPassword(hashedPassword, password)
func (m *MockPasswordHasher) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

// Mock createToken function to always success
func mockCreateToken(staff *pkg.Staff) (string, error) {
	return "mockTokenString", nil
}

func TestStaffService_CreateStaff(t *testing.T) {
	mockRepo := new(mockStaffRepo)
	service := staff.NewStaffService(mockRepo)

	// Test case: Successful staff creation
	t.Run("successful staff creation", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		mockRepo.On("CreateStaff", &inputStaff).Return(nil)

		createdStaff, err := service.CreateStaff(&inputStaff)

		assert.NoError(t, err)
		assert.Equal(t, "test_user", createdStaff.Username)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case: Failed - Create staff error
	t.Run("error staff repository creation", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		mockRepo.On("CreateStaff", &inputStaff).Return(errors.New("database error"))

		emptyStaff, err := service.CreateStaff(&inputStaff)

		assert.Error(t, err)
		assert.Empty(t, emptyStaff)

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestStaffService_SignInStaff(t *testing.T) {
	mockRepo := new(mockStaffRepo)
	mockHasher := new(MockPasswordHasher)

	// Create service instance
	service := &staff.StaffService{
		Repo:            mockRepo,
		PasswordHasher:  mockHasher,
		CreateTokenFunc: mockCreateToken,
	}

	// Test case: Successful staff sign in
	t.Run("successful staff sign in", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "secure_password",
		}

		// Mock the GetStaffFromUsername to return a staff member
		mockRepo.On("GetStaffFromUsername", inputStaff.Username).Return(&pkg.Staff{
			Username: "test_user",
			Password: "hash_password_is_oifdaifdhgiajgheahjrephjhg",
		}, nil)
		// Mock bcrypt CompareHashAndPassword to return nil (successful password match)
		mockHasher.On("CompareHashAndPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(nil)

		token, err := service.SignInStaff(&inputStaff)

		assert.NoError(t, err)
		assert.NotEmpty(t, token) // token is returned

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	// Test case: Failed - Cannot find email in database
	t.Run("no email in database", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "secure_password",
		}

		// Reset expectations for this test case
		mockRepo.ExpectedCalls = nil

		mockRepo.On("GetStaffFromUsername", inputStaff.Username).Return(nil, gorm.ErrRecordNotFound)

		token, err := service.SignInStaff(&inputStaff)

		assert.Error(t, err)
		assert.Empty(t, token) // token is empty string
		assert.EqualError(t, err, staff.ErrUnauthorized.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case: Failed - GetStaffFromUsername error in database
	t.Run("error staff repository lookup", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "secure_password",
		}

		// Reset expectations for this test case
		mockRepo.ExpectedCalls = nil

		mockRepo.On("GetStaffFromUsername", inputStaff.Username).Return(nil, errors.New("database error"))

		token, err := service.SignInStaff(&inputStaff)

		assert.Error(t, err)
		assert.Empty(t, token) // token is empty string
		assert.EqualError(t, err, "database error")

		// Verify expectations
		mockRepo.AssertExpectations(t)
	})

	// Test case: Failed - password is not matched
	t.Run("password is not matched", func(t *testing.T) {
		// mock input body request
		inputStaff := pkg.Staff{
			Username: "test_user",
			Password: "wrong_password",
		}

		// Reset expectations for this test case
		mockRepo.ExpectedCalls = nil
		mockHasher.ExpectedCalls = nil

		// Mock the GetStaffFromUsername to return a staff member
		mockRepo.On("GetStaffFromUsername", inputStaff.Username).Return(&pkg.Staff{
			Username: "test_user",
			Password: "wrong_hash_password_is_eataikugvkarugvtoy",
		}, nil)
		// Mock bcrypt CompareHashAndPassword to return error (failed password match)
		mockHasher.On("CompareHashAndPassword", mock.AnythingOfType("[]uint8"), mock.AnythingOfType("[]uint8")).Return(errors.New("password does not matched"))

		token, err := service.SignInStaff(&inputStaff)

		assert.Error(t, err)
		assert.Empty(t, token) // token is empty string
		assert.EqualError(t, err, staff.ErrUnauthorized.Error())

		// Verify expectations
		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})
}
