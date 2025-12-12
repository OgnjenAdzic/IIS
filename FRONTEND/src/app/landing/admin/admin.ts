import { Component, inject } from '@angular/core';
import { AuthService } from '../../auth/service/auth';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-admin',
  imports: [RouterLink],
  templateUrl: './admin.html',
  styleUrl: './admin.css',
})
export class Admin {
  authService = inject(AuthService);
}
