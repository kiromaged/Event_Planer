import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder } from '@angular/forms';
import { EventService, EventSearchFilters } from '../../Services/event.service';
import { TaskService } from '../../Services/task.service';
import { AuthService } from '../../Services/auth.service';

@Component({
  selector: 'app-search',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.css']
})
export class SearchComponent {
  results: { type: 'event' | 'task'; item: any }[] = [];
  isLoading = false;

  form: any;

  currentUserId: string | null = null;

  constructor(private fb: FormBuilder, private events: EventService, private tasks: TaskService, private auth: AuthService) {
    this.form = this.fb.group({
      keyword: [''],
      from: [''],
      to: [''],
      role: ['any']
    });
    this.currentUserId = this.auth.getCurrentUserId();
    this.auth.currentUser$.subscribe(id => (this.currentUserId = id));
  }

  search() {
    this.isLoading = true;
    const fv = this.form.value;
    const filters: EventSearchFilters = {
      keyword: fv.keyword || undefined,
      from: fv.from || undefined,
      to: fv.to || undefined,
      userId: this.currentUserId || undefined,
      role: (fv.role as EventSearchFilters['role']) || 'any'
    };

    // run both searches in sequence (could be parallel)
    this.events.searchEvents(filters).subscribe(evts => {
      this.tasks.searchTasks({ keyword: fv.keyword || undefined, from: fv.from || undefined, to: fv.to || undefined, userId: this.currentUserId || undefined, role: 'any' }).subscribe(ts => {
        this.results = [];
        evts.forEach(e => this.results.push({ type: 'event', item: e }));
        ts.forEach(t => this.results.push({ type: 'task', item: t }));
        this.isLoading = false;
      });
    });
  }
}
