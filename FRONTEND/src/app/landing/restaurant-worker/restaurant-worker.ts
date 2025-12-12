import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RestaurantService } from '../../features/restaurant/services/restaurant';
import { Restaurant } from '../../features/restaurant/models/restaurant';
import { RouterLink } from '@angular/router';

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
        console.log('Restaurants loaded:', this.restaurants);
        this.cdr.detectChanges();
      },
      error: (err) => {
        this.error = 'Failed to load restaurants';
        this.loading = false;
        console.error('Error loading restaurants:', err);
        this.cdr.detectChanges();
      }
    });
  }
}