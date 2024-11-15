package patient

import (
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"gorm.io/gorm"
)

// Secondary port
type PatientRepositoryInterface interface {
	SearchPatient() (*pkg.Patient, error)
}

// Secondary adapter
type GormPatientRepository struct {
	db *gorm.DB
}

// Initiate secondary adapter
func NewGormPatientRepository(db *gorm.DB) PatientRepositoryInterface {
	return &GormPatientRepository{db: db}
}

func (r *GormPatientRepository) SearchPatient() (*pkg.Patient, error) {
	//
}
