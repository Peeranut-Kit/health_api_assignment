package pkg

import (
	"time"
)

type Hospital struct {
	ID       int       `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"size:255" json:"name"`
	Patients []Patient `gorm:"foreignKey:HospitalID" json:"patients"`
	Staffs   []Staff   `gorm:"foreignKey:HospitalID" json:"staffs"`
}

type Patient struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	FirstNameTh  string     `gorm:"size:255" json:"first_name_th"`
	MiddleNameTh string     `gorm:"size:255" json:"middle_name_th"`
	LastNameTh   string     `gorm:"size:255" json:"last_name_th"`
	FirstNameEn  string     `gorm:"size:255" json:"first_name_en"`
	MiddleNameEn string     `gorm:"size:255" json:"middle_name_en"`
	LastNameEn   string     `gorm:"size:255" json:"last_name_en"`
	DateOfBirth  *time.Time `json:"date_of_birth"`
	PatientHN    string     `gorm:"size:50;not null;unique" json:"patient_hn"`
	NationalID   string     `gorm:"size:50;not null;unique" json:"national_id"`
	PassportID   string     `gorm:"size:50;not null;unique" json:"passport_id"`
	PhoneNumber  string     `gorm:"size:50;not null" json:"phone_number"`
	Email        string     `gorm:"size:255;unique" json:"email"`
	Gender       string     `gorm:"size:1" json:"gender"`
	HospitalID   int        `json:"hospital_id"`
	Hospital     Hospital   `gorm:"foreignKey:HospitalID" json:"hospital"`
}

type Staff struct {
	ID         int      `gorm:"primaryKey" json:"id"`
	Username   string   `gorm:"size:255;not null;unique" json:"username" validate:"required"`
	Password   string   `gorm:"size:255;not null" json:"password" validate:"required"`
	HospitalID int      `json:"hospital_id"`
	Hospital   Hospital `gorm:"foreignKey:HospitalID" json:"hospital"`
}
