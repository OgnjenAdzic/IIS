import { Component, inject } from '@angular/core';
import { AuthService } from '../../auth/service/auth';

@Component({
  selector: 'app-restaurant-worker',
  imports: [],
  templateUrl: './restaurant-worker.html',
  styleUrl: './restaurant-worker.css',
})
export class RestaurantWorker {
  authService = inject(AuthService);
}
