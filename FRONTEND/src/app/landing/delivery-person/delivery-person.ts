import { Component, inject } from '@angular/core';
import { AuthService } from '../../auth/service/auth';

@Component({
  selector: 'app-delivery-person',
  imports: [],
  templateUrl: './delivery-person.html',
  styleUrl: './delivery-person.css',
})
export class DeliveryPerson {
  authService = inject(AuthService);
}
