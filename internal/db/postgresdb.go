package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ApprovalRequest struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	ApprovalID  string    `gorm:"type:varchar(50);unique;not null"`
	Destination string    `gorm:"type:varchar(255);not null"`
	DateFrom    time.Time `gorm:"type:date;not null"`
	DateTo      time.Time `gorm:"type:date;not null"`
	Purpose     string    `gorm:"type:varchar(255);not null"`
	FlightInfo  string    `gorm:"type:text;default:''"`
	HotelInfo   string    `gorm:"type:text;default:''"`
	TotalCost   int       `gorm:"not null;default:0"`
	Notes       string    `gorm:"type:text;default:''"`
	Status      string    `gorm:"type:varchar(20);not null;default:'PENDING'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (ApprovalRequest) TableName() string {
	return "approval_requests"
}

func InitDB(dbUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the ApprovalRequest model
	if err := db.AutoMigrate(&ApprovalRequest{}); err != nil {
		return nil, err
	}

	return db, nil
}
