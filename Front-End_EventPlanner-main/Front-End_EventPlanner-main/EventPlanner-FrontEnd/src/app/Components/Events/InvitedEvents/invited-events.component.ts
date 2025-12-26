import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventService } from '../../../Services/event.service';
import { EventModel } from '../../../Models/event.model';
import { AuthService } from '../../../Services/auth.service';
import { EventItemComponent } from '../EventItem/event-item.component';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-invited-events',
  standalone: true,
  imports: [CommonModule, EventItemComponent],
  templateUrl: './invited-events.component.html',
  styleUrls: ['./invited-events.component.css']
})
export class InvitedEventsComponent implements OnInit, OnDestroy {
  events: EventModel[] = [];
  currentUserId: string | null = null;
  isLoading: boolean = false;
  errorMessage: string | null = null;
  private sub: Subscription | null = null;

  constructor(private svc: EventService, private auth: AuthService) {}

  ngOnInit(): void {
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
    this.svc.getInvitedEvents().subscribe({
      next: (list) => {
        this.events = list;
        this.isLoading = false;
      },
      error: (err) => {
        this.isLoading = false;
        if (err.status === 401) {
          this.errorMessage = 'You must be logged in to view your invited events.';
        } else {
          this.errorMessage = 'Failed to load invited events. Please try again.';
        }
        console.error('Load invited events error:', err);
      }
    });
  }

  setStatus(evtId: string, status: 'Going' | 'Maybe' | 'Not Going') {
    if (!this.currentUserId) {
      this.errorMessage = 'Not logged in';
      return;
    }

    // Convert frontend status to backend format
    const backendStatus = status === 'Going' ? 'going' : status === 'Maybe' ? 'maybe' : 'not_going';

    // Update UI immediately
    const event = this.events.find(e => e.id === evtId);
    if (event) {
      const attendeeIndex = event.attendees.findIndex(a => a.id === this.currentUserId);
      if (attendeeIndex >= 0) {
        event.attendees[attendeeIndex].status = status;
      } else {
        event.attendees.push({ id: this.currentUserId, status });
      }
    }

    // Then update on backend
    this.svc.setAttendanceStatus(evtId, backendStatus).subscribe({
      next: (ok) => {
        if (!ok) {
          this.errorMessage = 'Failed to set status';
          this.load();
        }
      },
      error: (err) => {
        console.error('Set status error:', err);
        this.errorMessage = 'Failed to update your attendance status.';
        this.load();
      }
    });
  }

  attendeeStatus(e: EventModel): string | undefined {
    if (!this.currentUserId) return undefined;
    return e.attendees.find(a => a.id === this.currentUserId)?.status;
  }

  isStatus(e: EventModel, status: 'Going' | 'Maybe' | 'Not Going') {
    return this.attendeeStatus(e) === status;
  }
}
