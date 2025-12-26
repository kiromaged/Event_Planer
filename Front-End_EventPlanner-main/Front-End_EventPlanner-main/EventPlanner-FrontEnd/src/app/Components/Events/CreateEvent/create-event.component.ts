import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { EventService } from '../../../Services/event.service';
import { AuthService } from '../../../Services/auth.service';

@Component({
  selector: 'app-create-event',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './create-event.component.html',
  styleUrls: ['./create-event.component.css']
})
export class CreateEventComponent {
  isLoading = false;
  errorMessage: string | null = null;
  form: any;
  currentUserId: string | null = null;

  constructor(private fb: FormBuilder, private svc: EventService, private router: Router, private auth: AuthService) {
    this.form = this.fb.group({
      title: ['', Validators.required],
      date: ['', Validators.required],
      time: ['', Validators.required],
      location: ['', Validators.required],
      description: ['']
    });
    // initialize current user from auth service
    this.currentUserId = this.auth.getCurrentUserId();
  }

  submit() {
    if (this.form.invalid) {
      this.errorMessage = 'Please fill in all required fields';
      return;
    }

    this.isLoading = true;
    this.errorMessage = null;
    const current = this.auth.getCurrentUserId();

    if (!current) {
      this.errorMessage = 'You must be logged in to create an event.';
      this.isLoading = false;
      return;
    }

    const formData = this.form.value;
    const data = {
      title: formData.title,
      date: formData.date,
      time: formData.time,
      location: formData.location,
      description: formData.description || '',
      organizerId: current,
      attendees: [] as string[]
    };

    // cast to required shape
    this.svc.create(data as any).subscribe({
      next: (evt) => {
        this.isLoading = false;
        // navigate to user's events
        this.router.navigate(['/events/mine']);
      },
      error: (err) => {
        this.isLoading = false;
        if (err.status === 400) {
          this.errorMessage = 'Invalid event data. Please check your input.';
        } else if (err.status === 401) {
          this.errorMessage = 'You must be logged in to create an event.';
        } else {
          this.errorMessage = 'Failed to create event. Please try again.';
        }
        console.error('Create event error:', err);
      }
    });
  }
}
