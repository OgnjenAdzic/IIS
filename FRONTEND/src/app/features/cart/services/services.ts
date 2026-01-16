import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import {
  CartItem,
  CartResponse,
  AddItemRequest,
  UpdateQuantityRequest
} from '../models/models';
import { MenuItem } from '../../restaurant/models/restaurant';

@Injectable({
  providedIn: 'root'
})
export class CartService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/cart`;

  // Signals remain for UI binding
  cartItems = signal<CartItem[]>([]);
  isOpen = signal<boolean>(false);
  cartRestaurantId = signal<string>('');

  // Total computed from backend data or frontend (for speed)
  totalPrice = computed(() => this.cartItems().reduce((acc, i) => acc + (i.price * i.quantity), 0));
  totalCount = computed(() => this.cartItems().reduce((acc, i) => acc + i.quantity, 0));

  constructor() {
    this.loadCart(); // Load on startup
  }

  toggleCart() {
    this.isOpen.update(v => !v);
  }

  // 1. LOAD FROM BACKEND
  loadCart() {
    this.http.get<CartResponse>(this.apiUrl).subscribe({
      next: (res) => {
        // Backend returns: { items: [...], totalPrice: ... }
        this.cartItems.set(res.items || []);
        this.cartRestaurantId.set(res.restaurantId || '');
      },
      error: (err) => console.error("Cart load failed", err)
    });
  }

  // 2. ADD ITEM (POST)
  addToCart(item: MenuItem, restaurantId: string) {
    const payload = {
      restaurantId: restaurantId,
      menuItemId: item.id,
      name: item.name,
      price: item.price,
      quantity: 1
    };

    this.http.post<CartResponse>(`${this.apiUrl}/items`, payload).subscribe({
      next: (res) => {
        this.cartItems.set(res.items); // Update UI with server response
        this.cartRestaurantId.set(res.restaurantId);
        this.isOpen.set(true); // Open Sidebar
      }
    });
  }

  // 3. REMOVE (DELETE)
  removeFromCart(itemId: string) {
    this.http.delete<CartResponse>(`${this.apiUrl}/items/${itemId}`).subscribe({
      next: (res) => {
        this.cartItems.set(res.items || []);

      }
    });
  }

  // 4. UPDATE QUANTITY (PUT)
  updateQuantity(itemId: string, delta: number) {
    const currentItem = this.cartItems().find(i => i.id === itemId);
    if (!currentItem) return;

    const newQty = currentItem.quantity + delta;
    if (newQty <= 0) {
      this.removeFromCart(itemId);
      return;
    }

    const payload = { itemId: itemId, quantity: newQty };
    this.http.put<CartResponse>(`${this.apiUrl}/items`, payload).subscribe({
      next: (res) => this.cartItems.set(res.items)
    });
  }

  clearCart() {
    this.http.delete<CartResponse>(this.apiUrl).subscribe({
      next: () => this.cartItems.set([])
    });
  }
}