import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { FormGroup, FormBuilder, Validators, ReactiveFormsModule } from '@angular/forms';
import { AuthService } from '../../Services/auth.service';

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css'],
  standalone: true,
  imports: [CommonModule, RouterModule, ReactiveFormsModule]
})
export class SignupComponent {
  signupForm: FormGroup;
  isLoading: boolean = false;
  errorMessage: string | null = null;

  constructor(private fb: FormBuilder, private auth: AuthService, private router: Router) {
    this.signupForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit() {
    if (this.signupForm.invalid) {
      this.errorMessage = 'Please fill in all required fields correctly';
      return;
    }

    this.isLoading = true;
    this.errorMessage = null;

    const { name, email, password } = this.signupForm.value;

    this.auth.signup(name, email, password).subscribe({
      next: (user) => {
        this.isLoading = false;
        // Auto-login after successful signup
        this.auth.login(email, password).subscribe({
          next: () => {
            this.router.navigate(['/events/mine']);
          },
          error: (err) => {
            console.error('Login after signup failed:', err);
            this.errorMessage = 'Signup successful, but please login manually.';
            this.router.navigate(['/login']);
          }
        });
      },
      error: (err) => {
        this.isLoading = false;
        if (err.status === 409) {
          this.errorMessage = 'Email already registered. Please use a different email.';
        } else if (err.status === 400) {
          this.errorMessage = 'Invalid input. Please check your details.';
        } else {
          this.errorMessage = 'Signup failed. Please try again.';
        }
        console.error('Signup error:', err);
      }
    });
  }
}