package models

import (
	"time"

	"gorm.io/datatypes"
)

// AuditLog represents an audit log entry in the PostgreSQL database
type AuditLog struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index:idx_audit_logs_user" json:"userId"`
	Action     string         `gorm:"not null;size:100" json:"action"`
	EntityType string         `gorm:"not null;size:50" json:"entityType"`
	EntityID   string         `gorm:"not null;size:255" json:"entityId"`
	Details    datatypes.JSON `gorm:"type:jsonb" json:"details"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index:idx_audit_logs_created" json:"createdAt"`
}

// TableName specifies the table name for the AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}
