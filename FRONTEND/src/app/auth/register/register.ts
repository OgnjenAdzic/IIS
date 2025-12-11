import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../service/auth';
import { UserRole } from '../model/auth';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './register.html',
  styleUrl: './register.css'
})
export class RegisterComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  // Make enum available in HTML
  roles = UserRole;
  errorMessage = '';

  registerForm = this.fb.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    role: [UserRole.CUSTOMER, Validators.required]
  });

  onSubmit() {
    if (this.registerForm.valid) {
      // Cast to any to satisfy the RegisterRequest interface structure
      const req = this.registerForm.value as any;

      this.authService.register(req).subscribe({
        next: () => {
          alert('Registration successful! Please login.');
          this.router.navigate(['/login']);
        },
        error: (err) => {
          console.error(err);
          this.errorMessage = 'Registration failed. Try a different username.';
        }
      });
    }
  }
}