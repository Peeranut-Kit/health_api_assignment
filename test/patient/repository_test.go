package patient

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Peeranut-Kit/health_api_assignment/internal/patient"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormPatientRepository_SearchPatient(t *testing.T) {
	// Mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	// GORM from mock database
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm database: %v", err)
	}

	repo := patient.NewGormPatientRepository(gormDB)

	// Success case
	t.Run("successful patient searching", func(t *testing.T) {
		// Mock input
		inputPatient := pkg.Patient{
			PatientHN:  "654350968",
			HospitalID: 1,
		}

		// Setup expectations
		mock.ExpectQuery("SELECT").WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "first_name_th", "middle_name_th", "last_name_th", "first_name_en", "middle_name_en", "last_name_en",
				"date_of_birth", "patient_hn", "national_id", "passport_id", "phone_number", "email", "gender", "hospital_id",
			}).
				AddRow(
					1, "พีรณัฐ", "กลาง", "กิตติวิทยากุล", "Peearnut", "Middle", "Kittivittayakul",
					time.Date(1997, 7, 31, 0, 0, 0, 0, time.UTC), "HN123456", "1234567890123", "P123456789", "0912345678", "max.pk@gmail.com", "M", 1,
				),
		)

		patientList, err := repo.SearchPatient(&inputPatient)

		assert.NoError(t, err) // err == nil
		assert.NotEmpty(t, patientList)
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Failure case
	t.Run("failed patient searching", func(t *testing.T) {
		// Mock input
		inputPatient := pkg.Patient{
			PatientHN:  "-7",
			HospitalID: 1,
		}

		// Setup expectations
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("database error"))

		patientList, err := repo.SearchPatient(&inputPatient)

		assert.Error(t, err)         // err happens and not nil
		assert.Empty(t, patientList) // patientList is empty
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
