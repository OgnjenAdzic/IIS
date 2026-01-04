import { Component, inject, OnInit } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from './navbar/navbar';
import { CartSidebarComponent } from './features/cart/components/cart-sidebar/cart-sidebar';
import { AuthService } from './auth/service/auth';
import { Stakeholders } from './features/stakeholders/service/stakeholders';

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
export class App implements OnInit {
  authService = inject(AuthService);
  stakeholdersService = inject(Stakeholders);

  ngOnInit() {
    if (this.authService.isLoggedIn()) {
      this.stakeholdersService.loadCurrentProfile();
    }
  }
}