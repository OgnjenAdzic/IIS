import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { AddItemForm } from '../../components/add-item-form/add-item-form';
import { MenuList } from '../../components/menu-list/menu-list';
import { StatusToggle } from '../../components/status-toggle/status-toggle';
import { RestaurantService } from '../../services/restaurant';
import { CommonModule } from '@angular/common';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-manage-restaurant',
  imports: [CommonModule, AddItemForm, MenuList, StatusToggle],
  templateUrl: './manage-restaurant.html',
  styleUrl: './manage-restaurant.css',
})
export class ManageRestaurant implements OnInit {
  private route = inject(ActivatedRoute);
  private restaurantService = inject(RestaurantService);
  private cdr = inject(ChangeDetectorRef);

  restaurant: any = null;
  restaurantId: string = '';

  ngOnInit() {
    // 1. Get the ID from the URL (e.g. /manage-restaurant/123)
    this.restaurantId = this.route.snapshot.paramMap.get('id') || '';

    if (this.restaurantId) {
      this.loadData();
    }
  }

  loadData() {
    this.restaurantService.getById(this.restaurantId).subscribe({
      next: (res) => {
        this.restaurant = res;
        console.log(res);
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error("Failed to load restaurant", err),
          this.cdr.detectChanges();
      }
    });
  }

  // --- EVENT HANDLERS (These were missing!) ---

  handleStatusChange(newStatus: boolean) {
    this.restaurantService.updateStatus(this.restaurantId, newStatus).subscribe(() => {
      // Update local state immediately
      this.restaurant.isOpen = newStatus;
      this.cdr.detectChanges();
    });
  }

  handleItemAdded(data: { name: string, price: number }) {
    this.restaurantService.addMenuItem(this.restaurantId, data.name, data.price)
      .subscribe(() => {
        this.loadData(); // Reload data to see the new item with its ID
      });
  }

  handlePriceUpdate(event: { id: string, price: number }) {
    this.restaurantService.updateItemPrice(event.id, event.price).subscribe(() => {
      console.log('Price updated successfully');
      this.cdr.detectChanges();
    });
  }

  handleItemDelete(itemId: string) {
    this.restaurantService.deleteMenuItem(itemId).subscribe(() => {
      // Remove item from the local array so we don't need to reload the whole page
      if (this.restaurant.menu && this.restaurant.menu.items) {
        this.restaurant.menu.items = this.restaurant.menu.items.filter((i: any) => i.id !== itemId);
        this.cdr.detectChanges();
      }
    });
  }

}
