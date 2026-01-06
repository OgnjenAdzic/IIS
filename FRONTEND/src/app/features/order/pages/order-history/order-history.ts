import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { OrderService } from '../../services/order.service';
import { Order } from '../../models/order.model';

@Component({
  selector: 'app-order-history',
  standalone: true,
  imports: [CommonModule, RouterLink, DatePipe],
  templateUrl: './order-history.html',
  styleUrl: './order-history.css'
})
export class OrderHistory implements OnInit {
  orderService = inject(OrderService);

  orders = signal<Order[]>([]);
  loading = signal<boolean>(true);

  private readonly activeStatuses = ['PENDING', 'PREPARING', 'READY_FOR_PICKUP', 'IN_DELIVERY'];

  sortedOrders = computed(() => {

    return [...this.orders()].sort((a, b) => {
      const aIsActive = this.activeStatuses.includes(a.status);
      const bIsActive = this.activeStatuses.includes(b.status);

      if (aIsActive && !bIsActive) return -1;
      if (!aIsActive && bIsActive) return 1;

      return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime();
    });
  });


  ngOnInit() {
    this.orderService.getMyOrders().subscribe({
      next: (res) => {
        this.orders.set(res.orders || []);
        this.loading.set(false);
      },
      error: (err) => {
        console.error(err);
        this.loading.set(false);
      }
    });
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'PENDING': return 'bg-yellow-100 text-yellow-800';
      case 'PREPARING': return 'bg-blue-100 text-blue-800';
      case 'IN_DELIVERY': return 'bg-purple-100 text-purple-800';
      case 'DELIVERED': return 'bg-green-100 text-green-800';
      case 'CANCELLED': return 'bg-red-100 text-red-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }
}