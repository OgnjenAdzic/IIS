import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterLink } from '@angular/router';
import { RestaurantService } from '../../features/restaurant/services/restaurant';

@Component({
  selector: 'app-customer',
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './customer.html',
  styleUrl: './customer.css',
})
export class Customer implements OnInit {
  restaurantService = inject(RestaurantService);
  restaurants = signal<any[]>([]);
  searchTerm = signal<string>('');
  onlyOpen = signal<boolean>(false);

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
    this.restaurantService.getAll().subscribe({
      next: (res: any) => {
        // Handle response structure { restaurants: [...] }
        const data = res.restaurants || res;
        this.restaurants.set(data);
      },
      error: (err) => console.error(err)
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
