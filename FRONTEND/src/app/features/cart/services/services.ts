import { Injectable, signal, computed } from '@angular/core';
import { CartItem } from '../models/models';

@Injectable({
  providedIn: 'root',
})
export class CartService {
  cartItems = signal<CartItem[]>([]);

  isOpen = signal<boolean>(false);

  totalCount = computed(() => {
    return this.cartItems().reduce((acc, item) => acc + item.quantity, 0);
  });

  totalPrice = computed(() => {
    return this.cartItems().reduce((acc, item) => acc + (item.price * item.quantity), 0);
  });


  toggleCart() {
    this.isOpen.update(val => !val);
  }

  addToCart(item: any, restaurantId: string) {
    this.cartItems.update(items => {
      const existing = items.find(i => i.id === item.id);

      if (existing) {
        return items.map(i =>
          i.id === item.id ? { ...i, quantity: i.quantity + 1 } : i
        );
      } else {
        return [...items, {
          id: item.id,
          name: item.name,
          price: item.price,
          quantity: 1,
          restaurantId: restaurantId
        }];
      }
    });

    this.isOpen.set(true);
  }

  removeFromCart(itemId: string) {
    this.cartItems.update(items => items.filter(i => i.id !== itemId));
  }

  updateQuantity(itemId: string, delta: number) {
    this.cartItems.update(items => {
      return items.map(item => {
        if (item.id === itemId) {
          const newQty = item.quantity + delta;
          return { ...item, quantity: newQty > 0 ? newQty : 1 };
        }
        return item;
      });
    });
  }

  clearCart() {
    this.cartItems.set([]);
  }

}
