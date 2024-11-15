package patient

import (
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
)

// Primary port
type PatientServiceInterface interface {
	SearchPatient(patientSearchRequest *pkg.Patient) ([]pkg.Patient, error)
}

type PatientService struct {
	repo PatientRepositoryInterface
}

func NewPatientService(repo PatientRepositoryInterface) PatientServiceInterface {
	return &PatientService{repo: repo}
}

func (s *PatientService) SearchPatient(patientSearchRequest *pkg.Patient) ([]pkg.Patient, error) {
	// Retrieve patient list searching
	patientList, err := s.repo.SearchPatient(patientSearchRequest)

	if err != nil {
		return nil, err
	}

	return patientList, nil
}
