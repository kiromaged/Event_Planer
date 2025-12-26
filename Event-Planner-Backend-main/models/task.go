package models

import "time"

// Task maps to the `tasks` table.
type Task struct {
	ID          uint       `gorm:"column:task_id;type:int unsigned;primaryKey;autoIncrement" json:"id"`
	EventID     uint       `gorm:"column:event_id;type:int unsigned;not null" json:"eventId"`
	Description string     `gorm:"column:description;type:text;not null" json:"description"`
	AssignedTo  *uint      `gorm:"column:assigned_to;type:int unsigned" json:"assignedTo,omitempty"`
	Status      string     `gorm:"column:status;type:enum('pending','in_progress','completed','cancelled');not null;default:'pending'" json:"status"`
	DueDate     *time.Time `gorm:"column:due_date;type:date" json:"dueDate,omitempty"`
	CreatedBy   uint       `gorm:"column:created_by;type:int unsigned;not null" json:"createdBy"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;autoUpdateTime" json:"updatedAt"`
	
	// Relations
	Event       Event      `gorm:"foreignKey:EventID" json:"event,omitempty"`
	Assignee    *User      `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Creator     User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName forces the GORM table name to `tasks`.
func (Task) TableName() string { return "tasks" }

