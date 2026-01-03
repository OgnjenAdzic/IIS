import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router'; // 1. Import Router

// Services
import { RestaurantService } from '../../services/restaurant'; // Fix path if needed
import { AuthService } from '../../../../auth/service/auth'; // 2. Import AuthService
import { UserRole } from '../../../../auth/model/auth';

// Components
import { AddItemForm } from '../../components/add-item-form/add-item-form';
import { MenuList } from '../../components/menu-list/menu-list';
import { StatusToggle } from '../../components/status-toggle/status-toggle';

@Component({
  selector: 'app-manage-restaurant',
  standalone: true,
  imports: [CommonModule, AddItemForm, MenuList, StatusToggle],
  templateUrl: './manage-restaurant.html',
  styleUrl: './manage-restaurant.css',
})
export class ManageRestaurant implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router); // 3. Inject Router
  private restaurantService = inject(RestaurantService);
  private authService = inject(AuthService); // 4. Inject AuthService
  private cdr = inject(ChangeDetectorRef);

  restaurant: any = null;
  restaurantId: string = '';

  ngOnInit() {
    this.restaurantId = this.route.snapshot.paramMap.get('id') || '';

    if (this.restaurantId) {
      this.loadData();
    }
  }

  loadData() {
    this.restaurantService.getById(this.restaurantId).subscribe({
      next: (res) => {
        // --- SECURITY CHECK ---
        const currentUser = this.authService.currentUser();

        // If user is NOT Admin AND User ID does NOT match Manager ID
        if (currentUser?.role !== UserRole.ADMIN && res.managerId !== currentUser?.id) {
          alert("Access Denied: You are not the manager of this restaurant.");
          this.router.navigate(['/restaurant']); // Redirect back to dashboard
          return;
        }
        // ----------------------

        this.restaurant = res;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error("Failed to load restaurant", err);
        // If restaurant doesn't exist, go back
        this.router.navigate(['/restaurant']);
      }
    });
  }

  // --- EVENT HANDLERS ---

  handleStatusChange(newStatus: boolean) {
    this.restaurantService.updateStatus(this.restaurantId, newStatus).subscribe({
      next: () => {
        this.restaurant.isOpen = newStatus;
        this.cdr.detectChanges();
      },
      error: (err) => alert("Error: You don't have permission to change status.")
    });
  }

  handleItemAdded(data: { name: string, price: number }) {
    this.restaurantService.addMenuItem(this.restaurantId, data.name, data.price)
      .subscribe({
        next: () => this.loadData(),
        error: (err) => alert("Error: You don't have permission to add items.")
      });
  }

  handlePriceUpdate(event: { id: string, price: number }) {
    this.restaurantService.updateItemPrice(event.id, event.price).subscribe({
      next: () => {
        console.log('Price updated successfully');
        this.cdr.detectChanges();
      },
      error: (err) => alert("Error updating price.")
    });
  }

  handleItemDelete(itemId: string) {
    this.restaurantService.deleteMenuItem(itemId).subscribe({
      next: () => {
        if (this.restaurant.menu && this.restaurant.menu.items) {
          this.restaurant.menu.items = this.restaurant.menu.items.filter((i: any) => i.id !== itemId);
          this.cdr.detectChanges();
        }
      },
      error: (err) => alert("Error deleting item.")
    });
  }
}