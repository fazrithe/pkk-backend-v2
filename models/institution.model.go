package models

import (
	"time"

	"github.com/google/uuid"
)

type Institution struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string    `gorm:"uniqueIndex;not null" json:"name,omitempty"`
	Address   string    `gorm:"not null" json:"address,omitempty"`
	Image     string    `json:"image,omitempty"`
	User      uuid.UUID `gorm:"not null" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreateInstitutionRequest struct {
	Name      string    `json:"name"  binding:"required"`
	Address   string    `json:"address" binding:"required"`
	Image     string    `json:"image"`
	User      string    `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type UpdateInstitution struct {
	Name      string    `json:"name,omitempty"`
	Address   string    `json:"address,omitempty"`
	Image     string    `json:"image,omitempty"`
	User      string    `json:"user,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type InstitutionResponse struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}
