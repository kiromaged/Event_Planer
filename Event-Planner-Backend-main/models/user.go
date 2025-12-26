package models

import "time"

// User maps exactly to the SQL `users` table defined in schema.sql.
//
// Table: users
// Columns:
// - user_id (INT UNSIGNED, PK, auto-increment)
// - name (VARCHAR(100), NOT NULL)
// - email (VARCHAR(255), NOT NULL, UNIQUE)
// - password_hash (VARCHAR(255), NOT NULL)
// - role (ENUM('organizer','attendee') NOT NULL DEFAULT 'attendee')
// - created_at (TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, indexed)
type User struct {
	ID           uint      `gorm:"column:user_id;type:int unsigned;primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Email        string    `gorm:"column:email;type:varchar(255);not null;uniqueIndex:ux_users_email" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"column:created_at;index:ix_users_created_at" json:"createdAt"`
}

// TableName forces the GORM table name to `users`.
func (User) TableName() string { return "users" }
