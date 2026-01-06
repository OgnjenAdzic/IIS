import { Component, Input, inject, OnInit, signal } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Order } from '../../../order/models/order.model';
import { OrderService } from '../../../order/services/order.service';

@Component({
  selector: 'app-incoming-orders',
  standalone: true,
  imports: [CommonModule, DatePipe, FormsModule],
  templateUrl: './incoming-orders.html',
  styleUrl: './incoming-orders.css'
})
export class IncomingOrders implements OnInit {
  @Input() restaurantId: string = '';
  private orderService = inject(OrderService);

  orders = signal<Order[]>([]);
  prepTimeInput: { [key: string]: number } = {}; // Store input for each order

  ngOnInit() {
    if (this.restaurantId) this.loadOrders();
  }

  loadOrders() {
    this.orderService.getRestaurantOrders(this.restaurantId).subscribe({
      next: (res) => { this.orders.set(res.orders || []), console.log(res.orders); },
      error: (err) => { console.error("Failed to load orders", err); }
    });
  }

  acceptOrder(orderId: string) {
    const minutes = this.prepTimeInput[orderId] || 15; // Default 15 mins
    this.orderService.updateStatus(orderId, 'PREPARING', minutes).subscribe(() => {
      this.loadOrders(); // Refresh list
    });
  }

  declineOrder(orderId: string) {
    if (confirm('Reject this order?')) {
      this.orderService.updateStatus(orderId, 'CANCELLED').subscribe(() => {
        this.loadOrders();
      });
    }
  }

  // Helper to mark ready for pickup
  markReady(orderId: string) {
    this.orderService.updateStatus(orderId, 'READY_FOR_PICKUP').subscribe(() => {
      this.loadOrders();
    });
  }
}