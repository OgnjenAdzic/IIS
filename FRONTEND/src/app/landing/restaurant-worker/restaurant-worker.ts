import { Component, inject, OnInit, ChangeDetectorRef, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RestaurantService } from '../../features/restaurant/services/restaurant';
import { Restaurant } from '../../features/restaurant/models/restaurant';
import { RouterLink } from '@angular/router';
import { OrderService } from '../../features/order/services/order.service';

@Component({
  selector: 'app-restaurant-worker',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './restaurant-worker.html',
  styleUrl: './restaurant-worker.css',
})
export class RestaurantWorker implements OnInit {
  private restaurantService = inject(RestaurantService);
  private cdr = inject(ChangeDetectorRef);
  private orderService = inject(OrderService);

  notificationCounts = signal<{ [key: string]: number }>({});

  restaurants: Restaurant[] = [];
  loading: boolean = true;
  error: string | null = null;

  ngOnInit(): void {
    this.loadRestaurants();
  }

  loadRestaurants(): void {
    this.restaurantService.getAll().subscribe({
      next: (response) => {
        this.restaurants = response.restaurants;
        this.loading = false;
        const data = response.restaurants || response;
        this.cdr.detectChanges();
        data.forEach((r: any) => this.checkPendingOrders(r.id));
      },
      error: (err) => {
        this.error = 'Failed to load restaurants';
        this.loading = false;
        console.error('Error loading restaurants:', err);
        this.cdr.detectChanges();
      }
    });
  }

  checkPendingOrders(restaurantId: string) {
    // Fetch only PENDING orders
    this.orderService.getRestaurantOrders(restaurantId, 'PENDING').subscribe({
      next: (res) => {
        const count = res.orders ? res.orders.length : 0;

        // Update the signal map
        this.notificationCounts.update(current => ({
          ...current,
          [restaurantId]: count
        }));
        this.cdr.detectChanges();
      }
    });
  }
}