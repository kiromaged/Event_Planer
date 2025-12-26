package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"event_planner_backend/config"
	"event_planner_backend/middleware"
	"event_planner_backend/models"
	"event_planner_backend/utils"
)

// SearchRequest represents the query parameters for advanced search.
type SearchRequest struct {
	Keyword  string `form:"keyword"`   // Search in event names and task descriptions
	Role     string `form:"role"`      // Filter by user role: "organizer" or "attendee"
	Type     string `form:"type"`      // "events" or "tasks" or "all" (default: "all")
}

// SearchEventsAndTasks performs advanced search on events and tasks.
func SearchEventsAndTasks(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, "invalid query parameters: "+err.Error())
		return
	}

	if config.DB == nil {
		utils.JSONError(c, http.StatusInternalServerError, "database not available")
		return
	}

	// Determine search type
	searchType := strings.ToLower(req.Type)
	if searchType == "" {
		searchType = "all"
	}

	response := gin.H{}

	// Search events if type is "events" or "all"
	if searchType == "events" || searchType == "all" {
		events, err := searchEvents(userID, req)
		if err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "failed to search events: "+err.Error())
			return
		}
		response["events"] = events
	}

	// Search tasks if type is "tasks" or "all"
	if searchType == "tasks" || searchType == "all" {
		tasks, err := searchTasks(userID, req)
		if err != nil {
			utils.JSONError(c, http.StatusInternalServerError, "failed to search tasks: "+err.Error())
			return
		}
		response["tasks"] = tasks
	}

	c.JSON(http.StatusOK, response)
}

// searchEvents searches events based on filters.
func searchEvents(userID uint, req SearchRequest) ([]gin.H, error) {
	query := config.DB.Model(&models.Event{}).
		Joins("INNER JOIN event_attendees ON events.event_id = event_attendees.event_id").
		Where("event_attendees.user_id = ?", userID)

	// Keyword search (title and description)
	if req.Keyword != "" {
		keyword := "%" + strings.ToLower(req.Keyword) + "%"
		query = query.Where("LOWER(events.title) LIKE ? OR LOWER(events.description) LIKE ?", keyword, keyword)
	}

	// Role filter
	if req.Role != "" {
		role := strings.ToLower(req.Role)
		if role == "organizer" || role == "attendee" {
			query = query.Where("event_attendees.role = ?", role)
		}
	}

	var events []models.Event
	if err := query.
		Preload("Organizer").
		Preload("Attendees", func(db *gorm.DB) *gorm.DB {
			return db.Preload("User")
		}).
		Select("events.*").
		Group("events.event_id").
		Order("events.created_at DESC").
		Find(&events).Error; err != nil {
		return nil, err
	}

	// Format response
	result := make([]gin.H, len(events))
	for i, event := range events {
		result[i] = formatEventResponse(event)
		// Find user's role and status from preloaded attendees
		for _, att := range event.Attendees {
			if att.UserID == userID {
				result[i]["myRole"] = att.Role
				result[i]["myStatus"] = att.Status
				break
			}
		}
		if event.CreatedBy == userID && result[i]["myRole"] == nil {
			result[i]["myRole"] = "organizer"
			result[i]["myStatus"] = "going"
		}
	}

	return result, nil
}

// searchTasks searches tasks based on filters.
func searchTasks(userID uint, req SearchRequest) ([]gin.H, error) {
	// Start with tasks that belong to events the user is part of
	query := config.DB.Model(&models.Task{}).
		Joins("INNER JOIN events ON tasks.event_id = events.event_id").
		Joins("INNER JOIN event_attendees ON events.event_id = event_attendees.event_id").
		Where("event_attendees.user_id = ?", userID)

	// Keyword search (task description)
	if req.Keyword != "" {
		keyword := "%" + strings.ToLower(req.Keyword) + "%"
		query = query.Where("LOWER(tasks.description) LIKE ?", keyword)
	}

	// Role filter (based on user's role in the event)
	if req.Role != "" {
		role := strings.ToLower(req.Role)
		if role == "organizer" || role == "attendee" {
			query = query.Where("event_attendees.role = ?", role)
		}
	}

	var tasks []models.Task
	if err := query.
		Preload("Event", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Organizer")
		}).
		Preload("Assignee").
		Preload("Creator").
		Select("tasks.*").
		Group("tasks.task_id").
		Order("tasks.due_date ASC, tasks.created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	// Format response
	result := make([]gin.H, len(tasks))
	for i, task := range tasks {
		dueDateStr := ""
		if task.DueDate != nil {
			dueDateStr = task.DueDate.Format("2006-01-02")
		}

		assigneeInfo := gin.H{}
		if task.AssignedTo != nil && task.Assignee != nil {
			assigneeInfo = gin.H{
				"id":    task.Assignee.ID,
				"name":  task.Assignee.Name,
				"email": task.Assignee.Email,
			}
		}

		result[i] = gin.H{
			"id":          task.ID,
			"eventId":     task.EventID,
			"eventTitle":  task.Event.Title,
			"description": task.Description,
			"assignedTo":  task.AssignedTo,
			"assignee":    assigneeInfo,
			"status":      task.Status,
			"dueDate":     dueDateStr,
			"createdBy":   task.CreatedBy,
			"creator": gin.H{
				"id":    task.Creator.ID,
				"name":  task.Creator.Name,
				"email": task.Creator.Email,
			},
			"createdAt": task.CreatedAt,
			"updatedAt": task.UpdatedAt,
		}
	}

	return result, nil
}

