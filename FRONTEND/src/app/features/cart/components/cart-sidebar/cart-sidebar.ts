import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CartService } from '../../services/services';

@Component({
  selector: 'app-cart-sidebar',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './cart-sidebar.html',
  styleUrl: './cart-sidebar.css'
})
export class CartSidebarComponent {
  cartService = inject(CartService);

  items = this.cartService.cartItems;
  total = this.cartService.totalPrice;
  isOpen = this.cartService.isOpen;

  close() {
    this.cartService.isOpen.set(false);
  }

  checkout() {
    alert("Proceeding to checkout with total: " + this.total());
  }
}