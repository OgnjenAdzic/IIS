import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { OrderService } from '../../services/order.service';
import { Order } from '../../models/order.model';

@Component({
  selector: 'app-order-details',
  standalone: true,
  imports: [CommonModule, DatePipe, RouterLink],
  templateUrl: './order-details.html',
  styleUrl: './order-details.css'
})
export class OrderDetails implements OnInit {
  private route = inject(ActivatedRoute);
  private orderService = inject(OrderService);

  order = signal<Order | null>(null);

  // Status steps for the UI
  steps = ['PENDING', 'PREPARING', 'READY_FOR_PICKUP', 'IN_DELIVERY', 'DELIVERED'];

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      // NOTE: In a real app, you should have a getOrderById endpoint.
      // For now, we reuse getMyOrders and find it locally.
      this.orderService.getMyOrders().subscribe(res => {
        const found = res.orders?.find(o => o.id === id);
        if (found) this.order.set(found);
      });
    }
  }

  isStepCompleted(step: string): boolean {
    const currentStatus = this.order()?.status || '';
    const currentIndex = this.steps.indexOf(currentStatus);
    const stepIndex = this.steps.indexOf(step);
    return stepIndex <= currentIndex;
  }

  isCurrentStep(step: string): boolean {
    return this.order()?.status === step;
  }
}