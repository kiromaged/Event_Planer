# EventPlanner Frontend - Complete Code Explanation

## 📚 Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture & Structure](#architecture--structure)
3. [Core Models](#core-models)
4. [Services](#services)
5. [Components](#components)
6. [Routing](#routing)
7. [Key Concepts](#key-concepts)

---

## Project Overview

**EventPlanner** is an Angular-based frontend application that allows users to:
- **Create events** and manage them
- **Invite users** to events they organize
- **View invited events** and respond with attendance status (Going, Maybe, Not Going)
- **Search events** by keywords and date ranges
- **Authenticate** via login/signup

### Tech Stack
- **Framework**: Angular 20.3.0 (latest version)
- **Language**: TypeScript 5.9
- **State Management**: RxJS (Observables)
- **Styling**: CSS
- **Forms**: Reactive Forms (FormBuilder)
- **Testing**: Jasmine & Karma

---

## Architecture & Structure

```
src/app/
├── Components/          # All UI components
│   ├── Events/         # Event-related components
│   │   ├── CreateEvent/
│   │   ├── MyEvents/
│   │   ├── InvitedEvents/
│   │   └── EventItem/
│   ├── login/
│   ├── signup/
│   ├── Search/
│   └── shared/
├── Services/           # Business logic & API calls
│   ├── auth.service.ts
│   ├── event.service.ts
│   ├── task.service.ts
│   └── users.service.ts
├── Models/             # TypeScript interfaces & types
│   ├── event.model.ts
│   └── task.model.ts
├── app.ts             # Root component
├── app.routes.ts      # Route configuration
└── app.config.ts      # Angular configuration
```

### Design Pattern: **Standalone Components + Services**
- Each component is **standalone** (modern Angular approach)
- **Services** handle all business logic and HTTP communication
- **Observables** manage state reactively

---

## Core Models

### EventModel
Located in: `src/app/Models/event.model.ts`

```typescript
export type AttendanceStatus = 'Going' | 'Maybe' | 'Not Going';

export interface Attendee {
  id: string;                    // User ID
  status: AttendanceStatus;      // User's RSVP status
}

export interface EventModel {
  id: string;                    // Unique event ID (generated)
  title: string;                 // Event name
  date: string;                  // ISO format: yyyy-mm-dd
  time?: string;                 // Optional: HH:MM format
  location?: string;             // Optional: event location
  description?: string;          // Optional: event details
  organizerId: string;           // ID of the user who created it
  attendees: Attendee[];         // List of people invited
}
```

**Key Points:**
- `organizerId` identifies who owns/created the event
- `attendees` is an array tracking who's invited and their response
- `id` is auto-generated when event is created

---

## Services

### 1. AuthService
**File**: `src/app/Services/auth.service.ts`
**Purpose**: Handle user authentication and user state

#### Key Methods:

```typescript
signup(name: string, email: string, password: string): Observable<UserPayload>
```
- Sends registration request to backend (`POST /api/signup`)
- Returns the newly created user

```typescript
login(email: string, password: string): Observable<{ token: string; user: UserPayload }>
```
- Sends login request to backend (`POST /api/login`)
- Stores token and user in localStorage
- Updates internal observables so components react to login

```typescript
logout(): void
```
- Clears token and user from localStorage
- Sets observables to null

```typescript
getCurrentUserId(): string | null
```
- Synchronous getter for current user ID
- Returns null if not logged in

#### Key Observables:

```typescript
currentUser$: Observable<string | null>
```
- Observable that emits the current user's ID
- Components subscribe to know when login/logout happens

```typescript
getCurrentUserProfile(): Observable<UserPayload | null>
```
- Returns full user object (id, name, email, role)

---

### 2. EventService
**File**: `src/app/Services/event.service.ts`
**Purpose**: Manage events (create, read, delete, invite)

#### Key Methods:

```typescript
create(data: Omit<EventModel, 'id'>): Observable<EventModel>
```
- Creates a new event
- Generates unique ID: `evt-XXXXXX`
- Adds to internal events array
- Notifies subscribers via BehaviorSubject

**Example usage:**
```typescript
this.svc.create({
  title: 'Birthday Party',
  date: '2025-12-25',
  time: '18:00',
  location: 'My House',
  description: 'Come celebrate!',
  organizerId: 'user-1',
  attendees: []
}).subscribe(event => console.log('Created:', event));
```

```typescript
getOrganizedEvents(userId: string): Observable<EventModel[]>
```
- Returns all events created by the given user
- Used in "My Events" component

```typescript
getInvitedEvents(userId: string): Observable<EventModel[]>
```
- Returns all events where user is invited (appears in attendees)
- Used in "Invited Events" component

```typescript
invite(eventId: string, userId: string): Observable<boolean>
```
- Adds a user to an event's attendees array
- Sets their initial status to 'Maybe'
- Returns true if successful

```typescript
setAttendanceStatus(eventId: string, userId: string, status: AttendanceStatus): Observable<boolean>
```
- Updates a user's RSVP status (Going/Maybe/Not Going)
- Creates attendee entry if user wasn't invited before

```typescript
delete(eventId: string, callerId: string): Observable<boolean>
```
- Only the organizer can delete
- Returns false if caller is not the organizer
- Removes event from array

```typescript
searchEvents(filters: EventSearchFilters = {}): Observable<EventModel[]>
```
- Advanced search with filters:
  - `keyword`: search in title/description
  - `from`: start date (ISO format)
  - `to`: end date (ISO format)
  - `userId`: filter by user (with role: 'organizer', 'attendee', 'any')

**Example:**
```typescript
this.svc.searchEvents({
  keyword: 'party',
  from: '2025-12-01',
  to: '2025-12-31',
  userId: 'user-1',
  role: 'organizer'  // Only events I organize
}).subscribe(events => console.log(events));
```

#### Internal State:
```typescript
private events: EventModel[] = [];  // In-memory event storage
private events$ = new BehaviorSubject<EventModel[]>(this.events);
```

**Important**: The service uses **in-memory storage** (not a real backend). In production, replace HTTP calls with real API endpoints.

---

## Components

### 1. Root Component (App)
**File**: `src/app/app.ts`

```typescript
@Component({
  selector: 'app-root',
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('EventPlanner-FrontEnd');
  
  constructor(public auth: AuthService) {}
}
```

**Purpose**: 
- Main app container
- Injects AuthService to make it available globally
- `<router-outlet>` renders the current page
- Navigation links for login, signup, etc.

**Signal**: `signal()` is Angular's new reactive primitive (like useState in React)

---

### 2. CreateEventComponent
**File**: `src/app/Components/Events/CreateEvent/create-event.component.ts`

```typescript
export class CreateEventComponent {
  isLoading = false;
  form: any;  // Reactive form
  currentUserId: string | null = null;

  constructor(
    private fb: FormBuilder,      // Build forms
    private svc: EventService,    // Create event
    private router: Router,       // Navigate after submit
    private auth: AuthService     // Get current user
  ) {
    this.form = this.fb.group({
      title: ['', Validators.required],
      date: ['', Validators.required],
      time: [''],
      location: [''],
      description: ['']
    });
    this.currentUserId = this.auth.getCurrentUserId();
  }

  submit() {
    if (this.form.invalid) return;  // Don't submit if form is invalid
    
    this.isLoading = true;
    const current = this.auth.getCurrentUserId();
    
    if (!current) {
      alert('You must be logged in to create an event.');
      this.isLoading = false;
      return;
    }

    // Prepare event data
    const data = {
      ...this.form.value,          // title, date, time, location, description
      organizerId: current,         // Set current user as organizer
      attendees: []                 // Start with no attendees
    };

    // Call service to create event
    this.svc.create(data as any).subscribe({
      next: (evt) => {
        this.isLoading = false;
        this.router.navigate(['/events/mine']);  // Go to "My Events"
      },
      error: () => (this.isLoading = false)
    });
  }
}
```

**Flow:**
1. User fills form (title, date, time, location, description)
2. Clicks "Create Event"
3. `submit()` is called
4. Validates form
5. Creates event via EventService
6. On success → navigates to "My Events" page

---

### 3. MyEventsComponent
**File**: `src/app/Components/Events/MyEvents/my-events.component.ts`

```typescript
export class MyEventsComponent implements OnInit, OnDestroy {
  events: EventModel[] = [];
  currentUserId: string | null = null;
  private sub: Subscription | null = null;

  constructor(
    private svc: EventService,
    private auth: AuthService,
    public users: UsersService
  ) {}

  ngOnInit(): void {
    // Subscribe to auth changes
    this.sub = this.auth.currentUser$.subscribe(uid => {
      this.currentUserId = uid;
      this.load();  // Reload events when user changes
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();  // Cleanup to prevent memory leaks
  }

  load() {
    if (!this.currentUserId) {
      this.events = [];
      return;
    }
    // Get all events organized by current user
    this.svc.getOrganizedEvents(this.currentUserId).subscribe(
      list => (this.events = list)
    );
  }

  deleteEvent(evt: EventModel) {
    if (!this.currentUserId) return alert('Not logged in');
    if (!confirm('Delete event "' + evt.title + '"?')) return;
    
    this.svc.delete(evt.id, this.currentUserId).subscribe(ok => {
      if (ok) {
        this.load();  // Reload list after deletion
      } else {
        alert('Unable to delete (not organizer)');
      }
    });
  }

  invitePrompt(evt: EventModel) {
    const userId = prompt('User ID to invite (example: user-2)');
    if (!userId) return;
    
    this.svc.invite(evt.id, userId).subscribe(ok => {
      if (ok) {
        alert('Invited ' + userId);
      } else {
        alert('Invite failed');
      }
    });
  }
}
```

**Key Lifecycle Hooks:**

- `ngOnInit()`: Called when component is created
  - Subscribe to current user changes
  - Load events whenever user logs in/out
  
- `ngOnDestroy()`: Called when component is destroyed
  - Unsubscribe from observables
  - Prevents memory leaks

**Features:**
- Shows only events YOU organized
- Can delete your events
- Can invite users to your events
- Auto-updates when you log in/out

---

### 4. InvitedEventsComponent
**File**: `src/app/Components/Events/InvitedEvents/invited-events.component.ts`

Similar to MyEventsComponent but:
- Shows events where YOU are invited (not organized)
- Uses `getInvitedEvents()` instead of `getOrganizedEvents()`
- Can respond with attendance status (Going/Maybe/Not Going)
- Cannot delete (only organizer can)

---

### 5. LoginComponent & SignupComponent
**Files**: 
- `src/app/Components/login/login.component.ts`
- `src/app/Components/signup/signup.component.ts`

**Purpose**: Authenticate users

**LoginComponent:**
- Form with email & password
- Calls `auth.login()`
- On success → navigate to home/events page
- On error → show error message

**SignupComponent:**
- Form with name, email & password
- Calls `auth.signup()`
- Creates new user account
- Typically redirects to login after signup

---

### 6. SearchComponent
**File**: `src/app/Components/Search/search.component.ts`

**Purpose**: Search for events with filters

**Features:**
- Keyword search (title/description)
- Date range filter (from/to)
- Filter by role (events you organize vs. attend)
- Results update as user types/changes filters

---

### 7. EventItemComponent
**File**: `src/app/Components/Events/EventItem/event-item.component.ts`

**Purpose**: Display a single event in a list

**Features:**
- Shows event title, date, time, location
- Shows organizer name
- Shows attendees and their statuses
- Buttons for actions (delete, invite, etc.)

---

## Routing

**File**: `src/app/app.routes.ts`

```typescript
export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  { path: 'events/create', component: CreateEventComponent },
  { path: 'events/mine', component: MyEventsComponent },
  { path: 'events/invited', component: InvitedEventsComponent },
  { path: 'search', component: SearchComponent }
];
```

**How it works:**
- Home (`/`) → redirects to `/login`
- `/login` → login form
- `/signup` → signup form
- `/events/create` → form to create new event
- `/events/mine` → list of events you organized
- `/events/invited` → list of events you're invited to
- `/search` → search for events

**Navigation in code:**
```typescript
this.router.navigate(['/events/mine']);  // Go to My Events
```

---

## Key Concepts

### 1. Observables & RxJS
**What**: Streams of data over time

**Example:**
```typescript
// Subscribe to auth changes
this.auth.currentUser$.subscribe(userId => {
  console.log('User changed to:', userId);
  // Automatically called when user logs in/out
});
```

**Why**: Lets UI automatically update when data changes

### 2. BehaviorSubject
**What**: An Observable that holds a value and emits new values

**Example:**
```typescript
private events$ = new BehaviorSubject<EventModel[]>([]);

// Get current value synchronously
const currentEvents = this.events$.getValue();

// Get as observable
this.events$.asObservable().subscribe(events => {
  console.log('Events changed:', events);
});
```

### 3. Reactive Forms (FormBuilder)
**What**: Type-safe way to create forms

**Example:**
```typescript
this.form = this.fb.group({
  title: ['', Validators.required],  // Default: '', Required
  date: ['', Validators.required],
  time: [''],  // Optional
});

// Access form value
const data = this.form.value;  // { title: '...', date: '...', time: '...' }

// Check if valid
if (this.form.valid) { ... }
```

### 4. Standalone Components
**What**: Components that don't need NgModule

**Example:**
```typescript
@Component({
  selector: 'app-my-component',
  standalone: true,  // ← No NgModule needed!
  imports: [CommonModule, ReactiveFormsModule],  // Import what you need
  template: '...',
  styles: '...'
})
export class MyComponent { }
```

### 5. Dependency Injection
**What**: Angular automatically provides services to components

**Example:**
```typescript
constructor(
  private eventService: EventService,  // ← Angular provides this
  private authService: AuthService
) { }
```

### 6. Unsubscribing (Memory Leak Prevention)
**What**: Clean up subscriptions to prevent memory leaks

**Bad:**
```typescript
ngOnInit() {
  this.auth.currentUser$.subscribe(...);  // ← Never unsubscribe = memory leak
}
```

**Good:**
```typescript
private sub: Subscription | null = null;

ngOnInit() {
  this.sub = this.auth.currentUser$.subscribe(...);
}

ngOnDestroy() {
  this.sub?.unsubscribe();  // ← Clean up!
}
```

---

## Data Flow Example: Creating an Event

```
User fills form
        ↓
User clicks "Create Event"
        ↓
CreateEventComponent.submit() called
        ↓
Validates form
        ↓
EventService.create() called
        ↓
Service creates unique ID
        ↓
Service adds to events array
        ↓
Service notifies subscribers (events$ Observable)
        ↓
Component receives response
        ↓
Router navigates to /events/mine
        ↓
MyEventsComponent loads
        ↓
Gets current userId from AuthService
        ↓
Calls EventService.getOrganizedEvents(userId)
        ↓
Service filters events by organizerId
        ↓
Component receives filtered list
        ↓
Template renders event list
```

---

## Data Flow Example: User Login

```
User enters email & password
        ↓
LoginComponent.submit() called
        ↓
AuthService.login(email, password) called
        ↓
Service sends POST /api/login
        ↓
Backend returns { token, user }
        ↓
Service saves token to localStorage
        ↓
Service saves user to localStorage
        ↓
Service updates currentUserId$ Observable
        ↓
All subscribed components react
        ↓
MyEventsComponent.load() called (subscribed to currentUser$)
        ↓
Events reload with new userId
```

---

## Common Tasks

### Add a new component
```bash
ng generate component Components/MyNewComponent
```

### Add a new service
```bash
ng generate service Services/my-new-service
```

### Run the app
```bash
ng serve
# Open http://localhost:4200
```

### Build for production
```bash
ng build
# Output in dist/
```

---

## Summary

**This is an event planning application where:**
- Users authenticate (login/signup)
- Authenticated users can create events
- Event organizers can invite others
- Invited users can RSVP (Going/Maybe/Not Going)
- Users can search for events
- All data is managed through Services using RxJS Observables
- UI automatically updates when data changes

**Architecture:**
- **Components** = UI views
- **Services** = Business logic & data
- **Models** = TypeScript types
- **Routes** = Navigation
- **RxJS** = Reactive state management

---

**Now you understand the complete codebase! Ask if you have questions about any specific part.** 🚀
