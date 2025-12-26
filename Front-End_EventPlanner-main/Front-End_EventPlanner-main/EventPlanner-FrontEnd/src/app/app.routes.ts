import { Routes } from '@angular/router';
import { LoginComponent } from './Components/login/login.component';
import { SignupComponent } from './Components/signup/signup.component';
import { CreateEventComponent } from './Components/Events/CreateEvent/create-event.component';
import { MyEventsComponent } from './Components/Events/MyEvents/my-events.component';
import { InvitedEventsComponent } from './Components/Events/InvitedEvents/invited-events.component';
import { SearchComponent } from './Components/Search/search.component';

export const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  { path: 'events/create', component: CreateEventComponent },
  { path: 'events/mine', component: MyEventsComponent },
  { path: 'events/invited', component: InvitedEventsComponent }
  ,{ path: 'search', component: SearchComponent }
];
