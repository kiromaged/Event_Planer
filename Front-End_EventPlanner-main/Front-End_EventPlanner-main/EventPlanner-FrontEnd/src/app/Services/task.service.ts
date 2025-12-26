import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { TaskModel } from '../Models/task.model';

@Injectable({ providedIn: 'root' })
export class TaskService {
  private baseUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  // Search tasks using backend endpoint
  searchTasks(filters: { keyword?: string; from?: string; to?: string; userId?: string; role?: 'assignee' | 'creator' | 'any' } = {}): Observable<TaskModel[]> {
    // Build query string
    let queryParams = new URLSearchParams();
    if (filters.keyword) queryParams.append('keyword', filters.keyword);
    if (filters.role && filters.role !== 'any') queryParams.append('role', filters.role);
    queryParams.append('type', 'tasks');

    return this.http.get<any>(`${this.baseUrl}/search?${queryParams.toString()}`).pipe(
      map(res => {
        const tasks = res.tasks || [];
        return tasks.map((t: any) => this.mapBackendTaskToFrontend(t));
      }),
      catchError(err => {
        console.error('Failed to search tasks:', err);
        return of([]);
      })
    );
  }

  // Helper to map backend task format to frontend TaskModel
  private mapBackendTaskToFrontend(backendTask: any): TaskModel {
    return {
      id: String(backendTask.id),
      title: backendTask.description || '',
      description: backendTask.description || '',
      date: backendTask.dueDate || '',
      createdBy: String(backendTask.createdBy),
      assigneeId: backendTask.assignedTo ? String(backendTask.assignedTo) : undefined
    };
  }
}
