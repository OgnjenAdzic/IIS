import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from './navbar/navbar';
import { CartSidebarComponent } from './features/cart/components/cart-sidebar/cart-sidebar';
import { AuthService } from './auth/service/auth';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, NavbarComponent, CartSidebarComponent],
  template: `
    <!-- Show Navbar only if logged in -->
    @if (authService.isLoggedIn()) {
      <app-navbar></app-navbar>
      
      <!-- Cart Sidebar (Always loaded, hidden via CSS until toggled) -->
      <app-cart-sidebar></app-cart-sidebar>
    }

    <router-outlet></router-outlet>
  `
})
export class App {
  authService = inject(AuthService);
}