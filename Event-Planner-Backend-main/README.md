# Event Planner Backend (Go + Gin + MySQL)

A comprehensive backend API for an Event Planner application built with Go, Gin framework, and MySQL database.

## Features

- **MySQL Database Support**: Uses MySQL exclusively via GORM
- **Authentication**: JWT-based authentication with bcrypt password hashing
- **Event Management**: Create, view, update, and delete events
- **User Invitations**: Invite users to events with role-based access
- **Attendance Tracking**: Track user attendance status (going, maybe, not_going)
- **Search Functionality**: Advanced search for events and tasks
- **CORS Enabled**: Configured for `http://localhost:4200`

## Technology Stack

- **Framework**: Gin (Go web framework)
- **Database**: MySQL (via GORM)
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **Environment**: `.env` file configuration

## Folder Structure

```
.
├── config/          # Environment and database configuration
├── controllers/     # Request handlers (auth, events, search, etc.)
├── models/          # GORM database models
├── routes/          # API route definitions
├── middleware/      # Authentication middleware
├── utils/           # Helper functions (JWT, password, responses)
├── event_planer_DB/ # Database schema SQL file
└── main.go          # Application entry point
```

## Prerequisites

- **Go 1.21+** installed on your system
- **MySQL 5.7+** or **MySQL 8.0+** database server
- **Postman** (optional, for API testing)

## Setup Instructions

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Database Setup

1. Create a MySQL database:
   ```sql
   CREATE DATABASE EventPlanner;
   ```

2. (Optional) Run the schema file to set up tables:
   ```bash
   mysql -u your_user -p EventPlanner < event_planer_DB/event_planer_schema.sql
   ```
   Or import it manually using MySQL Workbench or phpMyAdmin.

### 3. Environment Configuration

1. Copy the example environment file:
   ```bash
   cp env.example .env
   ```

2. Edit `.env` and configure your database:
   ```env
   DB_USER=your_mysql_username
   DB_PASSWORD=your_mysql_password
   DB_HOST=127.0.0.1
   DB_PORT=3306
   DB_NAME=EventPlanner
   JWT_SECRET=your_strong_secret_key_here
   ```

### 4. Run the Server

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

**Note**: If database environment variables are not set, the app will log a warning but continue running. However, most features require a database connection.

## API Endpoints

Base URL: `http://localhost:8080/api`

### Public Endpoints

#### Health Check
- **GET** `/api/ping`
  - **Description**: Verify server is running
  - **Response**: `{"message":"pong"}`

#### Signup
- **POST** `/api/signup`
  - **Description**: Create a new user account
  - **Request Body**:
    ```json
    {
      "name": "John Doe",
      "email": "john.doe@example.com",
      "password": "password123"
    }
    ```
  - **Success Response** (201):
    ```json
    {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@example.com"
    }
    ```
  - **Error Responses**:
    - `400`: Invalid payload
    - `409`: Email already registered
    - `500`: Server error

#### Login
- **POST** `/api/login`
  - **Description**: Authenticate and receive JWT token
  - **Request Body**:
    ```json
    {
      "email": "john.doe@example.com",
      "password": "password123"
    }
    ```
  - **Success Response** (200):
    ```json
    {
      "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john.doe@example.com"
      }
    }
    ```
  - **Error Responses**:
    - `400`: Invalid payload
    - `401`: Invalid email or password
    - `500`: Server error

### Protected Endpoints (Require Authentication)

All protected endpoints require the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

#### Events

##### Create Event
- **POST** `/api/events`
  - **Description**: Create a new event (user becomes organizer)
  - **Request Body**:
    ```json
    {
      "title": "Summer Music Festival",
      "description": "Annual summer music festival with multiple artists",
      "location": "Central Park, New York",
      "eventDate": "2024-07-15",
      "eventTime": "18:00:00"
    }
    ```
  - **Date Format**: `eventDate` must be in `YYYY-MM-DD` format (e.g., "2024-07-15")
  - **Time Format**: `eventTime` must be in `HH:MM:SS` or `HH:MM` format (e.g., "18:00:00" or "18:00")
  - **Success Response** (201):
    ```json
    {
      "id": 1,
      "title": "Summer Music Festival",
      "description": "Annual summer music festival",
      "location": "Central Park, New York",
      "eventDate": "2024-07-15",
      "eventTime": "18:00:00",
      "createdBy": 1,
      "createdAt": "2024-01-15T10:30:00Z"
    }
    ```

##### Get My Organized Events
- **GET** `/api/events/organized`
  - **Description**: Get all events created by the authenticated user
  - **Success Response** (200): Array of event objects

##### Get My Invited Events
- **GET** `/api/events/invited`
  - **Description**: Get all events the user is invited to
  - **Success Response** (200): Array of event objects with role and status

##### Get Event Details
- **GET** `/api/events/:id`
  - **Description**: Get detailed information about a specific event
  - **Success Response** (200): Event object with attendees list

##### Delete Event
- **DELETE** `/api/events/:id`
  - **Description**: Delete an event (only creator can delete)
  - **Success Response** (200):
    ```json
    {
      "message": "event deleted successfully"
    }
    ```

##### Invite User to Event
- **POST** `/api/events/:id/invite`
  - **Description**: Invite a user to an event (organizer only)
  - **Request Body**:
    ```json
    {
      "email": "jane.doe@example.com",
      "role": "attendee"
    }
    ```
  - **Role Options**: `"organizer"` or `"attendee"` (default: `"attendee"`)

#### Attendance

##### Update Attendance Status
- **PUT** `/api/events/:id/attendance`
  - **Description**: Update your attendance status for an event
  - **Request Body**:
    ```json
    {
      "status": "going"
    }
    ```
  - **Status Options**: `"going"`, `"maybe"`, `"not_going"`

##### Get Event Attendees
- **GET** `/api/events/:id/attendees`
  - **Description**: Get list of all attendees for an event (organizer only)
  - **Success Response** (200):
    ```json
    {
      "eventId": 1,
      "eventTitle": "Summer Music Festival",
      "attendees": [
        {
          "userId": 1,
          "userName": "John Doe",
          "userEmail": "john.doe@example.com",
          "role": "organizer",
          "status": "going",
          "invitedAt": "2024-01-15T10:30:00Z"
        }
      ]
    }
    ```

#### Search

##### Search Events and Tasks
- **GET** `/api/search`
  - **Description**: Search events and tasks with filters
  - **Query Parameters**:
    - `keyword` (optional): Search keyword for event titles, descriptions, or task descriptions
    - `type` (optional): `"events"`, `"tasks"`, or `"all"` (default: `"all"`)
    - `role` (optional): Filter by user role: `"organizer"` or `"attendee"`
  - **Example**: `/api/search?keyword=music&type=all&role=organizer`
  - **Success Response** (200):
    ```json
    {
      "events": [...],
      "tasks": [...]
    }
    ```

## Date and Time Format

**Important**: The API uses specific date and time formats:

- **Event Date**: Must be in `YYYY-MM-DD` format (e.g., `"2024-07-15"`)
- **Event Time**: Must be in `HH:MM:SS` or `HH:MM` format (e.g., `"18:00:00"` or `"18:00"`)

The backend automatically handles timezone conversion using the MySQL connection's local timezone settings.

## Testing with Postman

### Import Postman Collection

1. Open Postman
2. Click **Import** button
3. Select the file: `Event_Planner_API.postman_collection.json`
4. The collection will be imported with all endpoints organized in folders

### Postman Environment Setup

1. Create a new environment in Postman (or use the default)
2. Add the following variables:
   - `base_url`: `http://localhost:8080`
   - `auth_token`: (leave empty, will be auto-filled after login)
   - `event_id`: (leave empty, will be auto-filled after creating an event)
   - `user_id`: (leave empty, will be auto-filled after signup/login)

### Testing Walkthrough

#### Step 1: Health Check
1. Run the **Ping** request
2. Expected: `200 OK` with `{"message":"pong"}`

#### Step 2: User Registration
1. Run the **Signup** request
2. Update the email in the request body to a unique email
3. Expected: `201 Created` with user details
4. The `user_id` and `user_email` variables will be automatically set

#### Step 3: User Login
1. Run the **Login** request with the same credentials
2. Expected: `200 OK` with JWT token and user details
3. The `auth_token` variable will be automatically set for subsequent requests

#### Step 4: Create an Event
1. Run the **Create Event** request
2. **Important**: Ensure `eventDate` is in `YYYY-MM-DD` format (e.g., `"2024-07-15"`)
3. **Important**: Ensure `eventTime` is in `HH:MM:SS` or `HH:MM` format (e.g., `"18:00:00"`)
4. Expected: `201 Created` with event details
5. The `event_id` variable will be automatically set

#### Step 5: View Your Events
1. Run **Get My Organized Events** to see events you created
2. Run **Get My Invited Events** to see events you're invited to

#### Step 6: Invite Another User
1. First, create a second user account (signup with different email)
2. Login as the second user to get their token
3. Switch back to the first user's token
4. Run **Invite User to Event** with the second user's email
5. Expected: `201 Created` with invitation details

#### Step 7: Update Attendance Status
1. Login as the invited user (second user)
2. Run **Update Attendance Status** with status: `"going"`, `"maybe"`, or `"not_going"`
3. Expected: `200 OK` with updated status

#### Step 8: View Attendees
1. Login as the event organizer (first user)
2. Run **Get Event Attendees**
3. Expected: `200 OK` with list of all attendees

#### Step 9: Search
1. Run **Search Events and Tasks**
2. Try different query parameters:
   - `?keyword=music` - Search for "music" in events and tasks
   - `?type=events&keyword=festival` - Search only events
   - `?role=organizer` - Filter by organizer role
3. Expected: `200 OK` with search results

#### Step 10: Delete Event
1. As the event creator, run **Delete Event**
2. Expected: `200 OK` with success message

### Common Testing Scenarios

#### Test Date Format Validation
- Try creating an event with invalid date format (e.g., `"07/15/2024"`)
- Expected: `400 Bad Request` with error message about date format

#### Test Time Format Validation
- Try creating an event with invalid time format (e.g., `"6 PM"`)
- Expected: `400 Bad Request` with error message about time format

#### Test Authentication
- Try accessing a protected endpoint without the `Authorization` header
- Expected: `401 Unauthorized`

#### Test Authorization
- Try deleting an event you didn't create
- Expected: `403 Forbidden`

#### Test Duplicate Email
- Try signing up with an email that already exists
- Expected: `409 Conflict`

## Database Schema

The application uses the following main tables:

- **users**: User accounts
- **events**: Event information with `event_date` (DATE) and `event_time` (TIME)
- **event_attendees**: User-event relationships with roles and attendance status
- **tasks**: Tasks associated with events

See `event_planer_DB/event_planer_schema.sql` for the complete schema.

## Notes

- **MySQL Only**: This backend is configured exclusively for MySQL. The database driver is `gorm.io/driver/mysql`.
- **Password Security**: Passwords are hashed using bcrypt before storage.
- **JWT Tokens**: Tokens are signed using HS256 with the `JWT_SECRET` from environment variables. Tokens expire after 24 hours.
- **CORS**: Configured to allow requests from `http://localhost:4200` (Angular dev server).
- **Auto-migration**: On startup, GORM automatically creates/migrates tables based on the models.
- **Date Handling**: Event dates are stored as DATE type in MySQL and formatted as `YYYY-MM-DD` in API responses. Times are stored as TIME type and returned as strings.

## Troubleshooting

### Database Connection Issues
- Verify MySQL is running: `mysql -u root -p`
- Check `.env` file has correct credentials
- Ensure database `EventPlanner` exists
- Check firewall settings if connecting to remote MySQL

### Date/Time Format Errors
- Always use `YYYY-MM-DD` for dates (e.g., `"2024-07-15"`)
- Always use `HH:MM:SS` or `HH:MM` for times (e.g., `"18:00:00"` or `"18:00"`)
- Avoid using formats like `"07/15/2024"` or `"6 PM"`

### Authentication Errors
- Ensure JWT token is included in `Authorization: Bearer <token>` header
- Check if token has expired (tokens last 24 hours)
- Verify `JWT_SECRET` in `.env` matches the secret used when token was created

### Port Already in Use
- Change the port in `main.go` if `8080` is already in use
- Update `base_url` in Postman environment accordingly

## License

This project is provided as-is for educational and development purposes.
