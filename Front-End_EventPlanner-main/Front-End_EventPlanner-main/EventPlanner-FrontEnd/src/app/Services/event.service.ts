import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { EventModel } from '../Models/event.model';

export interface EventSearchFilters {
  keyword?: string; // search in title and description
  from?: string; // ISO date
  to?: string; // ISO date
  userId?: string; // user id for role-based filtering
  role?: 'organizer' | 'attendee' | 'any';
}

@Injectable({ providedIn: 'root' })
export class EventService {
  private baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  // Create a new event; returns observable of created event
  create(data: Omit<EventModel, 'id'>): Observable<EventModel> {
    // Map frontend format to backend format
    const payload = {
      title: data.title,
      description: data.description || '',
      location: data.location || '',
      eventDate: data.date, // Backend expects eventDate
      eventTime: data.time || '00:00' // Backend expects eventTime
    };
    return this.http.post<any>(`${this.baseUrl}/events`, payload).pipe(
      map(res => this.mapBackendEventToFrontend(res)),
      catchError(err => {
        console.error('Failed to create event:', err);
        return of();
      })
    );
  }

  // Get all events organized by the authenticated user
  getOrganizedEvents(): Observable<EventModel[]> {
    return this.http.get<any[]>(`${this.baseUrl}/events/organized`).pipe(
      map(res => res.map(e => this.mapBackendEventToFrontend(e))),
      catchError(err => {
        console.error('Failed to fetch organized events:', err);
        return of([]);
      })
    );
  }

  // Get all events where user is an attendee
  getInvitedEvents(): Observable<EventModel[]> {
    return this.http.get<any[]>(`${this.baseUrl}/events/invited`).pipe(
      map(res => res.map(e => this.mapBackendEventToFrontend(e))),
      catchError(err => {
        console.error('Failed to fetch invited events:', err);
        return of([]);
      })
    );
  }

  // Get event details by ID
  getEventDetails(eventId: string | number): Observable<EventModel> {
    return this.http.get<any>(`${this.baseUrl}/events/${eventId}`).pipe(
      map(res => this.mapBackendEventToFrontend(res)),
      catchError(err => {
        console.error('Failed to fetch event details:', err);
        return of();
      })
    );
  }

  // Invite a user to an event by email
  invite(eventId: string | number, email: string): Observable<boolean> {
    return this.http.post<any>(`${this.baseUrl}/events/${eventId}/invite`, { email, role: 'attendee' }).pipe(
      map(() => true),
      catchError(err => {
        console.error('Failed to invite user:', err);
        return of(false);
      })
    );
  }

  // Set attendance status for a user for the given event
  setAttendanceStatus(eventId: string | number, status: 'going' | 'maybe' | 'not_going'): Observable<boolean> {
    return this.http.put<any>(`${this.baseUrl}/events/${eventId}/attendance`, { status }).pipe(
      map(() => true),
      catchError(err => {
        console.error('Failed to update attendance status:', err);
        return of(false);
      })
    );
  }

  // Delete an event
  delete(eventId: string | number): Observable<boolean> {
    return this.http.delete<any>(`${this.baseUrl}/events/${eventId}`).pipe(
      map(() => true),
      catchError(err => {
        console.error('Failed to delete event:', err);
        return of(false);
      })
    );
  }

  // Get event attendees (organizer only)
  getEventAttendees(eventId: string | number): Observable<any> {
    return this.http.get<any>(`${this.baseUrl}/events/${eventId}/attendees`).pipe(
      catchError(err => {
        console.error('Failed to fetch attendees:', err);
        return of(null);
      })
    );
  }

  // Advanced search API for events
  searchEvents(filters: EventSearchFilters = {}): Observable<EventModel[]> {
    // Build query string
    let queryParams = new URLSearchParams();
    if (filters.keyword) queryParams.append('keyword', filters.keyword);
    if (filters.role && filters.role !== 'any') queryParams.append('role', filters.role);
    queryParams.append('type', 'events');

    return this.http.get<any>(`${this.baseUrl}/search?${queryParams.toString()}`).pipe(
      map(res => {
        const events = res.events || [];
        return events.map((e: any) => this.mapBackendEventToFrontend(e));
      }),
      catchError(err => {
        console.error('Failed to search events:', err);
        return of([]);
      })
    );
  }

  // Helper to map backend event format to frontend EventModel
  private mapBackendEventToFrontend(backendEvent: any): EventModel {
    return {
      id: String(backendEvent.id),
      title: backendEvent.title,
      date: backendEvent.eventDate, // Backend uses eventDate
      time: backendEvent.eventTime, // Backend uses eventTime
      location: backendEvent.location,
      description: backendEvent.description,
      organizerId: String(backendEvent.createdBy),
      attendees: (backendEvent.attendees || []).map((att: any) => ({
        id: String(att.userId),
        status: this.mapAttendanceStatus(att.status)
      }))
    };
  }

  // Map backend attendance status to frontend format
  private mapAttendanceStatus(backendStatus: string): 'Going' | 'Maybe' | 'Not Going' {
    const statusMap: Record<string, 'Going' | 'Maybe' | 'Not Going'> = {
      'going': 'Going',
      'maybe': 'Maybe',
      'not_going': 'Not Going',
      'pending': 'Maybe'
    };
    return statusMap[backendStatus.toLowerCase()] || 'Maybe';
  }
}
