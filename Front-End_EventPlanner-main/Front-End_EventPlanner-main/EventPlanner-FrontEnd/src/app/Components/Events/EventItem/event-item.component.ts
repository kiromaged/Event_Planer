import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { EventModel } from '../../../Models/event.model';
import { UsersService } from '../../../Services/users.service';

@Component({
  selector: 'app-event-item',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './event-item.component.html',
  styleUrls: ['./event-item.component.css']
})
export class EventItemComponent {
  @Input() event!: EventModel;
  // allow nullable current user id so templates can pass auth values safely
  @Input() currentUserId: string | null = null;

  constructor(private users: UsersService) {}

  get role(): string {
    if (!this.currentUserId) return '';
    if (this.event.organizerId === this.currentUserId) return 'Organizer';
    if (this.event.attendees?.some(a => a.id === this.currentUserId)) return 'Attendee';
    return '';
  }

  get organizerName(): string {
    return this.users.getDisplayName(this.event.organizerId);
  }
}
