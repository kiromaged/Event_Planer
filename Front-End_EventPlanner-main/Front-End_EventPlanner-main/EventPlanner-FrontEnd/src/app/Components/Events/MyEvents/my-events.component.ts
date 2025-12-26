import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventService } from '../../../Services/event.service';
import { AuthService } from '../../../Services/auth.service';
import { EventModel } from '../../../Models/event.model';
import { EventItemComponent } from '../EventItem/event-item.component';
import { Subscription } from 'rxjs';
import { UsersService } from '../../../Services/users.service';

@Component({
  selector: 'app-my-events',
  standalone: true,
  imports: [CommonModule, EventItemComponent],
  templateUrl: './my-events.component.html',
  styleUrls: ['./my-events.component.css']
})
export class MyEventsComponent implements OnInit, OnDestroy {
  events: EventModel[] = [];
  currentUserId: string | null = null;
  isLoading: boolean = false;
  errorMessage: string | null = null;
  private sub: Subscription | null = null;

  constructor(private svc: EventService, private auth: AuthService, public users: UsersService) {}

  ngOnInit(): void {
    // subscribe to auth changes and reload when user changes
    this.sub = this.auth.currentUser$.subscribe(uid => {
      this.currentUserId = uid;
      this.load();
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  load() {
    if (!this.currentUserId) {
      this.events = [];
      return;
    }

    this.isLoading = true;
    this.errorMessage = null;
    this.svc.getOrganizedEvents().subscribe({
      next: (list) => {
        this.events = list;
        this.isLoading = false;
      },
      error: (err) => {
        this.isLoading = false;
        if (err.status === 401) {
          this.errorMessage = 'You must be logged in to view your events.';
        } else {
          this.errorMessage = 'Failed to load events. Please try again.';
        }
        console.error('Load events error:', err);
      }
    });
  }

  deleteEvent(evt: EventModel) {
    if (!this.currentUserId) {
      this.errorMessage = 'Not logged in';
      return;
    }
    if (!confirm('Delete event "' + evt.title + '"?')) return;

    // Remove from UI immediately
    const index = this.events.findIndex(e => e.id === evt.id);
    if (index >= 0) {
      this.events.splice(index, 1);
    }

    // Then delete on backend
    this.svc.delete(evt.id).subscribe({
      next: (ok) => {
        if (!ok) {
          this.errorMessage = 'Unable to delete (not organizer)';
          this.load();
        }
      },
      error: (err) => {
        console.error('Delete error:', err);
        this.errorMessage = 'Failed to delete event';
        this.load();
      }
    });
  }

  invitePrompt(evt: EventModel) {
    const email = prompt('Email address to invite (example: user@example.com)');
    if (!email) return;

    this.svc.invite(evt.id, email).subscribe({
      next: (ok) => {
        if (ok) alert('Invited ' + email);
        else this.errorMessage = 'Invite failed';
      },
      error: (err) => {
        console.error('Invite error:', err);
        if (err.status === 404) {
          this.errorMessage = 'User not found with that email.';
        } else if (err.status === 409) {
          this.errorMessage = 'User is already invited to this event.';
        } else {
          this.errorMessage = 'Failed to send invite. Please try again.';
        }
      }
    });
  }
}
