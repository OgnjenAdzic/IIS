import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { RestaurantService } from '../../features/restaurant/services/restaurant';
import { PricingService } from '../../features/pricing/services/pricing.service';
import { Stakeholders } from '../../features/stakeholders/service/stakeholders';
import { AuthService } from '../../auth/service/auth';

@Component({
  selector: 'app-customer',
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './customer.html',
  styleUrl: './customer.css',
})
export class Customer implements OnInit {
  restaurantService = inject(RestaurantService);
  pricingService = inject(PricingService);
  stakeholdersService = inject(Stakeholders);
  authService = inject(AuthService);

  restaurants = signal<any[]>([]);
  searchTerm = signal<string>('');
  onlyOpen = signal<boolean>(false);

  customerLocation = signal<{ lat: number, lon: number } | null>(null);

  private priceInterval: any;

  // 2. COMPUTED SIGNAL (The Magic)
  // This automatically recalculates whenever one of the signals above changes
  filteredRestaurants = computed(() => {
    const term = this.searchTerm().toLowerCase();
    const showOpen = this.onlyOpen();
    const all = this.restaurants();

    return all.filter(r => {
      // Filter by Name or Category
      const matchesSearch = r.name.toLowerCase().includes(term) ||
        r.category.toLowerCase().includes(term);

      // Filter by Open Status (only if checkbox is checked)
      const matchesStatus = showOpen ? r.isOpen : true;

      return matchesSearch && matchesStatus;
    });
  });

  ngOnInit() {
    this.loadCustomerLocation();
    this.priceInterval = setInterval(() => {
      if (this.customerLocation()) {
        this.calculatePricesForList(this.restaurants());
      }
    }, 3000);
  }

  ngOnDestroy() {
    if (this.priceInterval) {
      clearInterval(this.priceInterval);
    }
  }

  loadCustomerLocation() {
    const user = this.authService.currentUser();
    if (user) {
      this.stakeholdersService.getCustomerProfile(user.id).subscribe({
        next: (profile) => {
          this.customerLocation.set({
            lat: profile.latitude,
            lon: profile.longitude
          });
          // Only load restaurants after we have location (to calc price)
          this.loadRestaurants();
        },
        error: () => this.loadRestaurants() // Load anyway even if location fails
      });
    }
  }

  loadRestaurants() {
    this.restaurantService.getAll().subscribe({
      next: (res: any) => {
        let data = res.restaurants || res;

        // If we have customer location, calculate price for each restaurant
        if (this.customerLocation()) {
          this.calculatePricesForList(data);
        } else {
          this.restaurants.set(data);
        }
      }
    });
  }

  calculatePricesForList(restaurants: any[]) {
    const custLoc = this.customerLocation()!;

    const updatedList = restaurants.map(r => ({
      ...r,
      deliveryPrice: null, // Placeholder
      loadingPrice: true
    }));

    this.restaurants.set(updatedList);

    // Now fetch prices one by one (or you could create a bulk endpoint)
    updatedList.forEach((r, index) => {
      this.pricingService.calculatePrice({
        customerLat: custLoc.lat,
        customerLon: custLoc.lon,
        restaurantLat: r.latitude,
        restaurantLon: r.longitude,
        isPriority: false // Default
      }).subscribe(priceRes => {
        // Update the signal array at specific index
        this.restaurants.update(currentList => {
          const newList = [...currentList];
          const found = newList.find(item => item.id === r.id);
          if (found) {
            found.deliveryPrice = priceRes.finalPrice;
            found.distanceKm = priceRes.distanceKm;
            found.loadingPrice = false;
          }
          return newList;
        });
      });
    });
  }


  // Helpers to update signals from HTML (optional, but cleaner)
  updateSearch(text: string) {
    this.searchTerm.set(text);
  }

  toggleOpen(checked: boolean) {
    this.onlyOpen.set(checked);
  }

}
