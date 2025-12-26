export type AttendanceStatus = 'Going' | 'Maybe' | 'Not Going';

export interface Attendee {
  id: string;
  status: AttendanceStatus;
}

export interface EventModel {
  id: string;
  title: string;
  date: string; // ISO date (yyyy-mm-dd)
  time?: string; // HH:MM
  location?: string;
  description?: string;
  organizerId: string; // user id of creator
  attendees: Attendee[]; // list of attendee objects with status
}
