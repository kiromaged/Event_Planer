# EventPlannerFrontEnd

This project was generated using [Angular CLI](https://github.com/angular/angular-cli) version 20.3.8.

## Development server

To start a local development server, run:

```bash
ng serve
```

Once the server is running, open your browser and navigate to `http://localhost:4200/`. The application will automatically reload whenever you modify any of the source files.

## Code scaffolding

Angular CLI includes powerful code scaffolding tools. To generate a new component, run:

```bash
ng generate component component-name
```

For a complete list of available schematics (such as `components`, `directives`, or `pipes`), run:

```bash
ng generate --help
```

## Building

To build the project run:

```bash
ng build
```

This will compile your project and store the build artifacts in the `dist/` directory. By default, the production build optimizes your application for performance and speed.

## Running unit tests

To execute unit tests with the [Karma](https://karma-runner.github.io) test runner, use the following command:

```bash
ng test
```

## Running end-to-end tests

For end-to-end (e2e) testing, run:

```bash
ng e2e
```

Angular CLI does not come with an end-to-end testing framework by default. You can choose one that suits your needs.

## Additional Resources

For more information on using the Angular CLI, including detailed command references, visit the [Angular CLI Overview and Command Reference](https://angular.dev/tools/cli) page.

## Component Documentation

### Login Component

The login component provides a secure authentication interface with the following features:

#### Structure
1. **Container Layout**
   - `auth-wrapper`: Main container for centering
   - `auth-card`: Elevated card design with form content

2. **Form Features**
   - Reactive Forms implementation
   - Real-time validation
   - Loading states
   - Error messaging

#### Key Elements
1. **Email Field**
   ```html
   <input 
       type="email" 
       formControlName="email"
       placeholder="name@company.com"
   >
   ```
   - Required field validation
   - Email format checking
   - Real-time error feedback

2. **Password Field**
   ```html
   <input 
       type="password"
       formControlName="password"
       placeholder="••••••••"
   >
   ```
   - Minimum 6 characters
   - Secure input handling
   - Validation feedback

3. **Submit Button**
   - Disabled state handling
   - Loading state indication
   - Form validation integration

4. **Navigation**
   - Sign-up link for new users
   - Router integration
   - Clean state management

#### Implementation Details

The component uses Angular's Reactive Forms for robust form handling:

```typescript
loginForm = this.fb.group({
  email: ['', [Validators.required, Validators.email]],
  password: ['', [Validators.required, Validators.minLength(6)]]
});
```

Key validation features:
- Email format validation
- Password length requirements
- Required field handling
- Touch state tracking
- Loading state management

#### User Experience
- Clear error messages
- Visual feedback for invalid fields
- Disabled states for invalid forms
- Loading indication during submission
- Smooth navigation flow

#### Best Practices
1. **Accessibility**
   - Proper label associations
   - ARIA attributes
   - Semantic HTML structure

2. **Security**
   - Secure password handling
   - Form validation
   - Protected routes

3. **Performance**
   - Efficient form updates
   - Optimized validation
   - Smart state management

#### Routes
- `/login` - Login page
- `/signup` - Registration page
- `/` - Protected home page

## Project organization

I added a short `PROJECT_STRUCTURE.md` at the repository root that explains where to find models, services, and components. Key places:

- `src/app/Models` — data interfaces and `index.ts` barrel
- `src/app/Services` — services and `index.ts` barrel
- `src/app/Components` — components organized by feature with a top-level `index.ts` barrel

If you're reviewing the repo, start with `PROJECT_STRUCTURE.md` for a quick map of the codebase.

## Event Management (Organizer role)

This project includes a basic, in-memory implementation of event management for the Organizer role. The implementation is intended as a client-side scaffold you can replace with backend APIs later.

Files added:

- `src/app/Models/event.model.ts` — TypeScript interface describing event data (id, title, date, time, location, description, organizerId, attendees[]).
- `src/app/Services/event.service.ts` — In-memory service exposing create, getOrganizedEvents, getInvitedEvents, invite, delete and allEvents$.
- `src/app/Components/Events/create-event.component.*` — Form to create a new event (organizer becomes current user).
- `src/app/Components/Events/my-events.component.*` — Lists events organized by the current user and allows deleting and inviting.
- `src/app/Components/Events/invited-events.component.*` — Lists events where the current user is an attendee.
- `src/app/Components/Events/event-item.component.ts` — Small reusable display component showing organizer/attendee badge.

How it works (client-side stub):

1. Create: `CreateEventComponent` collects title, date, time, location and description. It calls `EventService.create()` with `organizerId` set to the current user id (in the stub `user-1`).
2. Organizer view: `MyEventsComponent` calls `EventService.getOrganizedEvents(currentUserId)` to list events the user created. The organizer may delete an event (if they are the creator) or invite another user by id.
3. Invited view: `InvitedEventsComponent` calls `EventService.getInvitedEvents(currentUserId)` to list events where the user was added to attendees.

Notes & next steps:
- Replace the in-memory `EventService` with an HTTP-backed service that calls your server API for persistent storage and authentication-aware operations.
- Integrate with your auth module to use the real `currentUserId` (instead of the hard-coded `user-1`).
- Add proper user selection UI when inviting (instead of simple prompt). Add email invites if needed.
- Implement Authorization on the server so only the organizer can delete an event.

