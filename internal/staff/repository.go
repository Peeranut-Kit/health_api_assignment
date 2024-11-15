package staff

import (
	"github.com/Peeranut-Kit/health_api_assignment/pkg"
	"gorm.io/gorm"
)

// Secondary port
type StaffRepositoryInterface interface {
	CreateStaff(staff *pkg.Staff) error
	GetStaffFromUsername(username string) (*pkg.Staff, error)
}

// Secondary adapter
type GormStaffRepository struct {
	db *gorm.DB
}

// Initiate secondary adapter
func NewGormStaffRepository(db *gorm.DB) StaffRepositoryInterface {
	return &GormStaffRepository{db: db}
}

func (r *GormStaffRepository) CreateStaff(staff *pkg.Staff) error {
	// Create to staff database
	if result := r.db.Create(&staff); result.Error != nil {
		return result.Error
	}

	/*// Preload the associated hospital
	if err := r.db.Preload("Hospital").First(&staff, staff.ID).Error; err != nil {
		return err
	}*/

	return nil
}

func (r *GormStaffRepository) GetStaffFromUsername(username string) (*pkg.Staff, error) {
	var staff pkg.Staff
	if err := r.db.Where("username = ?", username).First(&staff).Error; err != nil {
		return nil, err
	}

	return &staff, nil
}
