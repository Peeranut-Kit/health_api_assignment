package staff

import (
	"errors"
	"os"
	"time"

	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrUnauthorized = errors.New("unauthorization. Username or password is wrong")

// Additional interface to mock for testing
// PasswordHasher is an interface that wraps bcrypt.CompareHashAndPassword.
type PasswordHasher interface {
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}

// BcryptHasher is a concrete implementation of PasswordHasher that uses bcrypt.
type BcryptHasher struct{}

func (b *BcryptHasher) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

// Start of service original file

// Primary port
type StaffServiceInterface interface {
	CreateStaff(staff *pkg.Staff) (*pkg.Staff, error)
	SignInStaff(staff *pkg.Staff) (string, error)
}

type StaffService struct {
	Repo            StaffRepositoryInterface
	PasswordHasher  PasswordHasher
	CreateTokenFunc func(staff *pkg.Staff) (string, error)
}

func NewStaffService(repo StaffRepositoryInterface) StaffServiceInterface {
	return &StaffService{
		Repo:            repo,
		PasswordHasher:  &BcryptHasher{},
		CreateTokenFunc: createToken,
	}
}

func (s *StaffService) CreateStaff(staff *pkg.Staff) (*pkg.Staff, error) {
	// encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// re-assign user password before saving in database
	staff.Password = string(hashedPassword)

	if err := s.Repo.CreateStaff(staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *StaffService) SignInStaff(staff *pkg.Staff) (string, error) {
	// Retrieve user by email
	selectedStaffByEmail, err := s.Repo.GetStaffFromUsername(staff.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrUnauthorized
		} else {
			return "", err
		}
	}

	// Compare the provided password with the hash stored in the database
	if err := s.PasswordHasher.CompareHashAndPassword([]byte(selectedStaffByEmail.Password), []byte(staff.Password)); err != nil {
		return "", ErrUnauthorized
	}

	// Create JWT token for the authenticated staff
	token, err := s.CreateTokenFunc(selectedStaffByEmail)
	if err != nil {
		return "", errors.New("error creating token")
	}

	// Success Sign In
	return token, nil
}

func createToken(staff *pkg.Staff) (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["staff_id"] = staff.ID
	claims["staff_name"] = staff.Username
	claims["staff_hospital_id"] = staff.HospitalID
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expiration 1 hour from now

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it in response. (t is token)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret)) // <- Secret key (keep this safe!)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
