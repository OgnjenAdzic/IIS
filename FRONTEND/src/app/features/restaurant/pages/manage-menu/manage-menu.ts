import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';

import { RestaurantService } from '../../services/restaurant';
import { AuthService } from '../../../../auth/service/auth';
import { UserRole } from '../../../../auth/model/auth';

// Only Menu Components
import { AddItemForm } from '../../components/add-item-form/add-item-form';
import { MenuList } from '../../components/menu-list/menu-list';

@Component({
  selector: 'app-manage-menu',
  standalone: true,
  imports: [CommonModule, AddItemForm, MenuList, RouterLink],
  templateUrl: './manage-menu.html',
  styleUrl: './manage-menu.css'
})
export class ManageMenu implements OnInit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private restaurantService = inject(RestaurantService);
  private authService = inject(AuthService);
  private cdr = inject(ChangeDetectorRef);

  restaurant: any = null;
  restaurantId: string = '';

  ngOnInit() {
    this.restaurantId = this.route.snapshot.paramMap.get('id') || '';
    if (this.restaurantId) this.loadData();
  }

  loadData() {
    this.restaurantService.getById(this.restaurantId).subscribe({
      next: (res) => {
        // Same Security Check
        const currentUser = this.authService.currentUser();
        if (currentUser?.role !== UserRole.ADMIN && res.managerId !== currentUser?.id) {
          this.router.navigate(['/restaurant']);
          return;
        }
        this.restaurant = res;
        this.cdr.detectChanges();
      },
      error: () => this.router.navigate(['/restaurant'])
    });
  }

  handleItemAdded(data: { name: string, price: number }) {
    this.restaurantService.addMenuItem(this.restaurantId, data.name, data.price)
      .subscribe(() => this.loadData());
  }

  handlePriceUpdate(event: { id: string, price: number }) {
    this.restaurantService.updateItemPrice(event.id, event.price).subscribe();
  }

  handleItemDelete(itemId: string) {
    this.restaurantService.deleteMenuItem(itemId).subscribe(() => {
      if (this.restaurant?.menu?.items) {
        this.restaurant.menu.items = this.restaurant.menu.items.filter((i: any) => i.id !== itemId);
        this.cdr.detectChanges();
      }
    });
  }
}