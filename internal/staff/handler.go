package staff

import (
	"net/http"

	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/gin-gonic/gin"
)

// Primary adapter
type StaffHandler struct {
	service StaffServiceInterface
}

// Just define what struct will do
type StaffHandlerInterface interface {
	CreateStaff(c *gin.Context)
	SignInStaff(c *gin.Context)
}

func NewHttpStaffHandler(service StaffServiceInterface) *StaffHandler {
	return &StaffHandler{service: service}
}

func (h *StaffHandler) CreateStaff(c *gin.Context) {
	var newStaff pkg.Staff
	/*if err := c.BindJSON(&newStaff); err != nil {
	    // The error is automatically handled by Gin.
	    return
	}*/
	if err := c.ShouldBindJSON(&newStaff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	createdStaff, err := h.service.CreateStaff(&newStaff)

	// Internal service error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Success creation
	c.JSON(http.StatusCreated, gin.H{
		"message": "Created successfully",
		"data":    createdStaff,
	})
}

func (h *StaffHandler) SignInStaff(c *gin.Context) {
	var staffInput pkg.Staff
	if err := c.ShouldBindJSON(&staffInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call service
	token, err := h.service.SignInStaff(&staffInput)

	// Internal service error
	if err != nil {
		if err == ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Set the JWT token in a cookie
	c.SetCookie("jwt", token, 3600, "/", "localhost", false, true)

	// Success login response
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
