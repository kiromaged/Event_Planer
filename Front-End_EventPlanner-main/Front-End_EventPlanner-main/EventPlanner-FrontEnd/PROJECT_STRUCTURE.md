Project structure

This file describes how the repository is organized so that any developer or reviewer can quickly find code and understand responsibilities.

Top-level
- `angular.json`, `package.json`, TypeScript config files: standard Angular project files.
- `README.md`: high-level project info and running instructions.

Key source folders (`src`)
- `src/app/Models/` - TypeScript interfaces and types used across the app.
  - `event.model.ts` - Event data shape (includes `Attendee` and `AttendanceStatus`).
  - `task.model.ts` - Task data shape.
  - `index.ts` - Barrel file exporting models.

- `src/app/Services/` - Application services (in-memory stubs now).
  - `event.service.ts` - Event CRUD, invites, attendance status and search API.
  - `task.service.ts` - Task CRUD and search API.
  - `auth.service.ts` - Simple auth stub exposing `currentUser$`.
  - `users.service.ts` - Static user directory for display names.
  - `index.ts` - Barrel file exporting services.

- `src/app/Components/` - UI components arranged by feature.
  - `login/`, `signup/` - Authentication components.
  - `Events/` - Event-related components (create, list, invited, item view).
  - `Search/` - Global search component.
  - `index.ts` - Barrel exporting top-level components and sub-barrels.

App wiring
- `src/app/app.routes.ts` - Central route definitions (add new pages here).
- `src/app/app.ts` and `src/app/app.html` - Application shell, imports, and header nav.

How to navigate quickly
- To find event-related UI: `src/app/Components/Events/`
- To find event business logic and search: `src/app/Services/event.service.ts`
- To update data shapes: `src/app/Models/event.model.ts`

Development notes
- Services are in-memory stubs. Replace with HTTP-backed implementations in `src/app/Services/`.
- Components are standalone (Angular 20 style) to simplify tree-shaking and per-component imports.

Common tasks
- Start dev server: `npm start` (runs `ng serve`)
- Build production: `npm run build`

If you'd like, I can:
- Add a README section per component documenting props and usage, or
- Create an automated script to generate a file tree snapshot for easy review.

*** End of file
