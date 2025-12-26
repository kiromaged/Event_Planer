package main

import (
	"log"

	"event_planner_backend/config"
	"event_planner_backend/models"
	"event_planner_backend/routes"
)

// main starts the HTTP server.
func main() {
	config.LoadEnv()
	config.InitDB()

	// Auto-migrate if DB is connected; safe no-op otherwise
	if config.DB != nil {
		if err := config.DB.AutoMigrate(
			&models.User{},
			&models.Event{},
			&models.EventAttendee{},
			&models.Task{},
		); err != nil {
			log.Printf("auto-migrate failed: %v", err)
		}
	}

	r := routes.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
