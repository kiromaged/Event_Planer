package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"event_planner_backend/controllers"
	"event_planner_backend/middleware"
)

// SetupRouter configures routes and middleware.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS for Angular dev server and OpenShift
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://*.openshiftapps.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			// Allow localhost and any OpenShift route
			return origin == "http://localhost:4200" || 
				   len(origin) > 0 && (origin[:8] == "https://" || origin[:7] == "http://")
		},
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api")
	{
		// Public routes
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		api.POST("/signup", controllers.Signup)
		api.POST("/login", controllers.Login)

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Event routes
			protected.POST("/events", controllers.CreateEvent)
			protected.GET("/events/organized", controllers.GetMyOrganizedEvents)
			protected.GET("/events/invited", controllers.GetMyInvitedEvents)
			protected.GET("/events/:id", controllers.GetEventDetails)
			protected.DELETE("/events/:id", controllers.DeleteEvent)
			protected.POST("/events/:id/invite", controllers.InviteUserToEvent)

			// Response/Attendance routes
			protected.PUT("/events/:id/attendance", controllers.UpdateAttendanceStatus)
			protected.GET("/events/:id/attendees", controllers.GetEventAttendees)

			// Search routes
			protected.GET("/search", controllers.SearchEventsAndTasks)
		}
	}

	return r
}
