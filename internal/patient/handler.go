package patient

import (
	"github.com/gin-gonic/gin"
)

// Primary adapter
type PatientHandler struct {
	service PatientServiceInterface
}

// Just define what struct will do
type PatientHandlerInterface interface {
	SearchPatient(c *gin.Context)
}

func NewHttpPatientHandler(service PatientServiceInterface) *PatientHandler {
	return &PatientHandler{service: service}
}

func (h *PatientHandler) SearchPatient(c *gin.Context) {
	//
}
