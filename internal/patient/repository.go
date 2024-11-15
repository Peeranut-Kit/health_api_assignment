package patient

import (
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"gorm.io/gorm"
)

// Secondary port
type PatientRepositoryInterface interface {
	SearchPatient(request *pkg.Patient) ([]pkg.Patient, error)
}

// Secondary adapter
type GormPatientRepository struct {
	db *gorm.DB
}

// Initiate secondary adapter
func NewGormPatientRepository(db *gorm.DB) PatientRepositoryInterface {
	return &GormPatientRepository{db: db}
}

func (r *GormPatientRepository) SearchPatient(request *pkg.Patient) ([]pkg.Patient, error) {
	var patientList []pkg.Patient

	query := r.db.Table("patients").Where("hospital_id = ?", request.HospitalID)

	// Add optional conditions only if fields are populated
    if request.ID != 0 {
        query = query.Where("id = ?", request.ID)
    }
    if request.FirstNameTh != "" {
        query = query.Where("first_name_th = ?", request.FirstNameTh)
    }
    if request.MiddleNameTh != "" {
        query = query.Where("middle_name_th = ?", request.MiddleNameTh)
    }
    if request.LastNameTh != "" {
        query = query.Where("last_name_th = ?", request.LastNameTh)
    }
    if request.FirstNameEn != "" {
        query = query.Where("first_name_en = ?", request.FirstNameEn)
    }
    if request.MiddleNameEn != "" {
        query = query.Where("middle_name_en = ?", request.MiddleNameEn)
    }
    if request.LastNameEn != "" {
        query = query.Where("last_name_en = ?", request.LastNameEn)
    }
    if request.DateOfBirth != nil {
        query = query.Where("date_of_birth = ?", *request.DateOfBirth)
    }
    if request.PatientHN != "" {
        query = query.Where("patient_hn = ?", request.PatientHN)
    }
    if request.NationalID != "" {
        query = query.Where("national_id = ?", request.NationalID)
    }
    if request.PassportID != "" {
        query = query.Where("passport_id = ?", request.PassportID)
    }
    if request.PhoneNumber != "" {
        query = query.Where("phone_number = ?", request.PhoneNumber)
    }
    if request.Email != "" {
        query = query.Where("email = ?", request.Email)
    }
    if request.Gender != "" {
        query = query.Where("gender = ?", request.Gender)
    }

	// Execute the query
    if err := query.Find(&patientList).Error; err != nil {
        return nil, err
    }

    return patientList, nil
}
