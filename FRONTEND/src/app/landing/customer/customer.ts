import { Component, inject } from '@angular/core';
import { AuthService } from '../../auth/service/auth';

@Component({
  selector: 'app-customer',
  imports: [],
  templateUrl: './customer.html',
  styleUrl: './customer.css',
})
export class Customer {
  authService = inject(AuthService);

}
