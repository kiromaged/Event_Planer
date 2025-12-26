import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';

export interface UserPayload {
  id: number | string;
  name: string;
  email: string;
  role?: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  // Backend base URL
  private baseUrl = 'http://localhost:8080/api';
  private tokenKey = 'auth_token';
  private userKey = 'auth_user';

  // Profile subject (full user object)
  private currentUserProfile$ = new BehaviorSubject<UserPayload | null>(null);

  // Compatibility subject: many components subscribe to `currentUser$` expecting a string id
  private currentUserId$ = new BehaviorSubject<string | null>(null);

  // Keep the old public name `currentUser$` as the id Observable (string|null) for compatibility
  get currentUser$(): Observable<string | null> {
    return this.currentUserId$.asObservable();
  }

  constructor(private http: HttpClient) {
    // Restore user/profile from localStorage if present
    const raw = localStorage.getItem(this.userKey);
    if (raw) {
      try {
        const user: UserPayload = JSON.parse(raw);
        this.currentUserProfile$.next(user);
        this.currentUserId$.next(String(user.id));
      } catch (_) {
        localStorage.removeItem(this.userKey);
      }
    }
  }

  // Signup against backend: POST /api/signup -> returns created user (201)
  signup(name: string, email: string, password: string): Observable<UserPayload> {
    return this.http.post<UserPayload>(`${this.baseUrl}/signup`, { name, email, password });
  }

  // Login against backend: POST /api/login -> { token, user }
  login(email: string, password: string): Observable<{ token: string; user: UserPayload }> {
    return this.http.post<{ token: string; user: UserPayload }>(`${this.baseUrl}/login`, { email, password }).pipe(
      tap(res => {
        if (res?.token) {
          localStorage.setItem(this.tokenKey, res.token);
        }
        if (res?.user) {
          localStorage.setItem(this.userKey, JSON.stringify(res.user));
          this.currentUserProfile$.next(res.user);
          this.currentUserId$.next(String(res.user.id));
        }
      })
    );
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
    localStorage.removeItem(this.userKey);
    this.currentUserProfile$.next(null);
    this.currentUserId$.next(null);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  // Compatibility synchronous getter used in some components
  getCurrentUserId(): string | null {
    return this.currentUserId$.getValue();
  }

  // Expose full profile observable for parts of the app that need the user object
  getCurrentUserProfile(): Observable<UserPayload | null> {
    return this.currentUserProfile$.asObservable();
  }

  // Synchronous convenience getter for profile
  getCurrentUserProfileValue(): UserPayload | null {
    return this.currentUserProfile$.getValue();
  }
}
