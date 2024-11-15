package staff

import (
	"net/http"

	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// CreateStaff godoc
// @Summary Create a new staff member
// @Description Create a new hospital staff member with login credentials
// @Tags Staff
// @Accept json
// @Produce json
// @Param staff body pkg.Staff true "Staff details"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /staff/create [post]
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

	// Validate the input body
	validate := validator.New()
	err := validate.Struct(newStaff)
	if err != nil {
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

// SignInStaff godoc
// @Summary Staff login
// @Description Authenticates a staff member and returns a JWT token
// @Tags Staff
// @Accept json
// @Produce json
// @Param credentials body pkg.Staff true "Staff login credentials"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /staff/login [post]
func (h *StaffHandler) SignInStaff(c *gin.Context) {
	var staffInput pkg.Staff
	if err := c.ShouldBindJSON(&staffInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the input body
	validate := validator.New()
	err := validate.Struct(staffInput)
	if err != nil {
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
