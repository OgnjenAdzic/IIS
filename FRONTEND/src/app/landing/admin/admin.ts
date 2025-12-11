import { Component, inject } from '@angular/core';
import { AuthService } from '../../auth/service/auth';

@Component({
  selector: 'app-admin',
  imports: [],
  templateUrl: './admin.html',
  styleUrl: './admin.css',
})
export class Admin {
  authService = inject(AuthService);
}
