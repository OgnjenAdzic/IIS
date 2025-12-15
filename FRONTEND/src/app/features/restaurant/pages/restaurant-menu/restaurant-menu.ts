import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { RestaurantService } from '../../services/restaurant';
import { Restaurant } from '../../models/restaurant';
import { CartService } from '../../../cart/services/services';

@Component({
  selector: 'app-restaurant-menu',
  imports: [CommonModule, RouterLink],
  templateUrl: './restaurant-menu.html',
  styleUrl: './restaurant-menu.css',
})
export class RestaurantMenu implements OnInit {
  private route = inject(ActivatedRoute);
  private restaurantService = inject(RestaurantService);
  private cartService = inject(CartService);

  restaurant = signal<Restaurant | null>(null);
  loading = signal<boolean>(true);

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadRestaurant(id);
    }
  }

  loadRestaurant(id: string) {
    this.restaurantService.getById(id).subscribe({
      next: (res: Restaurant) => {
        this.restaurant.set(res);
        this.loading.set(false);
      },
      error: (err) => {
        console.error(err);
        this.loading.set(false);
      }
    });
  }

  // Placeholder for future Cart logic
  addToCart(item: any) {
    const currentRestaurant = this.restaurant();
    if (currentRestaurant && currentRestaurant.id) {
      this.cartService.addToCart(item, currentRestaurant.id);
    } else {
      console.error("Error: Restaurant data is not loaded.");
    }
  }

}
