import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AuthService } from '../auth/service/auth';
import { CartService } from '../features/cart/services/services';
import { UserRole } from '../auth/model/auth';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: './navbar.html',
  styleUrl: './navbar.css'
})
export class NavbarComponent {
  authService = inject(AuthService);
  cartService = inject(CartService);

  user = this.authService.currentUser;
  cartCount = this.cartService.totalCount;

  Roles = UserRole;

  getDashboardRoute(): string {
    const role = this.user()?.role;
    switch (role) {
      case this.Roles.CUSTOMER:
        return '/customer';
      case this.Roles.ADMIN:
        return '/admin';
      case this.Roles.DELIVERY_PERSON:
        return '/delivery';
      case this.Roles.RESTAURANT_WORKER:
        return '/restaurant';
      default:
        return '/login';
    }
  }

  toggleCart() {
    this.cartService.toggleCart();
  }

  logout() {
    this.authService.logout();
  }
}