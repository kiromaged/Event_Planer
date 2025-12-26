package models

import "time"

// EventAttendee maps to the `event_attendees` table.
// This table stores the relationship between users and events,
// including their role (organizer/attendee) and attendance status.
type EventAttendee struct {
	EventID    uint      `gorm:"column:event_id;type:int unsigned;primaryKey" json:"eventId"`
	UserID     uint      `gorm:"column:user_id;type:int unsigned;primaryKey" json:"userId"`
	Role       string    `gorm:"column:role;type:enum('organizer','attendee');not null;default:'attendee'" json:"role"`
	Status     string    `gorm:"column:status;type:enum('going','maybe','not_going','pending');not null;default:'pending'" json:"status"`
	InvitedAt  time.Time `gorm:"column:invited_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"invitedAt"`
	
	// Relations
	Event      Event     `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName forces the GORM table name to `event_attendees`.
func (EventAttendee) TableName() string { return "event_attendees" }

