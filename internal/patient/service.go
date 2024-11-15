package patient

import (
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
)

// Primary port
type PatientServiceInterface interface {
	SearchPatient() (*pkg.Patient, error)
}

type PatientService struct {
	repo PatientRepositoryInterface
}

func NewPatientService(repo PatientRepositoryInterface) PatientServiceInterface {
	return &PatientService{repo: repo}
}

func (s *PatientService) SearchPatient() (*pkg.Patient, error) {
	//
}
