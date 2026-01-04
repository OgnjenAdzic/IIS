import { Component, inject, OnDestroy, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { RestaurantService } from '../../services/restaurant';
import { Restaurant } from '../../models/restaurant';
import { CartService } from '../../../cart/services/services';
import { Stakeholders } from '../../../stakeholders/service/stakeholders';
import { PricingService } from '../../../pricing/services/pricing.service';

@Component({
  selector: 'app-restaurant-menu',
  imports: [CommonModule, RouterLink],
  templateUrl: './restaurant-menu.html',
  styleUrl: './restaurant-menu.css',
})
export class RestaurantMenu implements OnInit, OnDestroy {
  private route = inject(ActivatedRoute);
  private restaurantService = inject(RestaurantService);
  private cartService = inject(CartService);
  private stakeholdersService = inject(Stakeholders);
  private pricingService = inject(PricingService);

  restaurant = signal<Restaurant | null>(null);
  loading = signal<boolean>(true);

  deliveryInfo = signal<{ price: number; distance: number } | null>(null);
  calculatingPrice = signal<boolean>(false);

  private priceInterval: any;

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');

    // 1. Initial Load from Router State (Instant)
    const state = history.state;
    if (state.price && state.distance) {
      this.deliveryInfo.set({ price: state.price, distance: state.distance });
    }

    if (id) {
      this.loadRestaurant(id);
    }

    // 4. START INTERVAL (e.g., every 30 seconds)
    this.priceInterval = setInterval(() => {
      const currentRestaurant = this.restaurant();
      if (currentRestaurant) {
        console.log("Auto-refreshing delivery price...");
        this.calculateDeliveryCost(currentRestaurant);
      }
    }, 3000);
  }

  ngOnDestroy() {
    if (this.priceInterval) {
      clearInterval(this.priceInterval);
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

  addToCart(item: any) {
    const currentRestaurant = this.restaurant();
    if (currentRestaurant && currentRestaurant.id) {
      this.cartService.addToCart(item, currentRestaurant.id);
    } else {
      console.error("Error: Restaurant data is not loaded.");
    }
  }

  calculateDeliveryCost(restaurant: Restaurant) {
    const coords = this.stakeholdersService.getCoordinates();

    if (!coords) {
      console.warn("User location not found. Is profile loaded?");
      // Optional: Call loadCurrentProfile() here if missing
      return;
    }

    this.calculatingPrice.set(true);

    // 2. Call Pricing API directly
    this.pricingService.calculatePrice({
      customerLat: coords.lat,
      customerLon: coords.lon,
      restaurantLat: restaurant.latitude,
      restaurantLon: restaurant.longitude,
      isPriority: false
    }).subscribe({
      next: (priceRes) => {
        this.deliveryInfo.set({
          price: priceRes.finalPrice,
          distance: priceRes.distanceKm
        });
        this.calculatingPrice.set(false);
      },
      error: () => this.calculatingPrice.set(false)
    });
  }

}
