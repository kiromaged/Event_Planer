package models

import "time"

// Event maps to the `events` table.
type Event struct {
	ID          uint      `gorm:"column:event_id;type:int unsigned;primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Location    string    `gorm:"column:location;type:varchar(255);not null" json:"location"`
	EventDate   time.Time `gorm:"column:event_date;type:date;not null" json:"eventDate"`
	EventTime   string    `gorm:"column:event_time;type:varchar(8);not null" json:"eventTime"`
	CreatedBy   uint      `gorm:"column:created_by;type:int unsigned;not null" json:"createdBy"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`

	// Relations
	Organizer User            `gorm:"foreignKey:CreatedBy" json:"organizer,omitempty"`
	Attendees []EventAttendee `gorm:"foreignKey:EventID" json:"attendees,omitempty"`
}

// TableName forces the GORM table name to `events`.
func (Event) TableName() string { return "events" }
