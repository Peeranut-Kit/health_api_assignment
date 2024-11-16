package staff_test

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Peeranut-Kit/health_api_assignment/internal/staff"
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGormStaffRepository_CreateStaff(t *testing.T) {
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

	repo := staff.NewGormStaffRepository(gormDB)

	// Success case
	t.Run("successful staff creation", func(t *testing.T) {
		// Mock input
		newStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		// Setup expectations for the mock database
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "staffs"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Call the repository function
		err := repo.CreateStaff(&newStaff)

		// err == nil
		assert.NoError(t, err)
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Failure case
	t.Run("failed staff creation", func(t *testing.T) {
		// Mock input
		newStaff := pkg.Staff{
			Username:   "test_user",
			Password:   "secure_password",
			HospitalID: 1,
		}

		// Setup expectations
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "staffs"`).WillReturnError(errors.New("database error"))
		mock.ExpectRollback()

		err = repo.CreateStaff(&newStaff)

		// err happens and not nil
		assert.Error(t, err)
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGormStaffRepository_GetStaffFromUsername(t *testing.T) {
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

	repo := staff.NewGormStaffRepository(gormDB)

	// Success case
	t.Run("successful staff retrieving from username", func(t *testing.T) {
		// Mock input
		username := "test_username"

		// Setup expectations
		//mock.ExpectBegin() no need to ExpectBegin() because SELECT does not use transaction like INSERT
		mock.ExpectQuery("SELECT").WillReturnRows(
			sqlmock.NewRows(
				[]string{"id", "username", "password", "hospital_id"}).
				AddRow(1, "test_username", "hashed_password", 1))
		//mock.ExpectCommit()

		staff, err := repo.GetStaffFromUsername(username)

		assert.NoError(t, err)	// err == nil
		assert.NotEmpty(t, staff)	// staff not empty
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Failure case
	t.Run("failed staff retrieving from username", func(t *testing.T) {
		// Mock input
		username := "wrong_test_username"

		// Setup expectations
		mock.ExpectQuery("SELECT").WillReturnError(errors.New("database error"))

		staff, err := repo.GetStaffFromUsername(username)

		assert.Error(t, err)	// err happens and not nil
		assert.Empty(t, staff)	// staff is empty
		// Ensure all expectations were met
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
