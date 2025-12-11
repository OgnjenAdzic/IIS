import { Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../service/auth';
import { UserRole } from '../model/auth';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink], // We don't need CommonModule with new @if syntax
  templateUrl: './login.html',
  styleUrl: './login.css'
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private authService = inject(AuthService);
  private router = inject(Router);

  errorMessage = '';

  loginForm = this.fb.group({
    username: ['', Validators.required],
    password: ['', Validators.required]
  });

  onSubmit() {
    if (this.loginForm.valid) {
      const { username, password } = this.loginForm.value;

      this.authService.login({ username: username!, password: password! })
        .subscribe({
          next: () => {
            const user = this.authService.currentUser();
            this.redirectBasedOnRole(user?.role);
            console.log("uspesno");
            console.log(user);
            console.log(user?.role);
          },
          error: (err) => {
            console.error(err);
            this.errorMessage = 'Invalid username or password';
          }
        });
    }
  }

  private redirectBasedOnRole(role?: UserRole) {
    switch (role) {
      case UserRole.ADMIN:
        this.router.navigate(['/admin']);
        break;
      case UserRole.CUSTOMER:
        this.router.navigate(['/customer']);
        break;
      case UserRole.DELIVERY_PERSON:
        this.router.navigate(['/delivery']);
        break;
      case UserRole.RESTAURANT_WORKER:
        this.router.navigate(['/restaurant']);
        break;
      default:
        this.router.navigate(['/login']); // Fallback
    }
  }
}