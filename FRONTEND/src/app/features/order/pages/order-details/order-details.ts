import { Component, inject, OnInit, signal, OnDestroy } from '@angular/core';
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
export class OrderDetails implements OnInit, OnDestroy {
  private route = inject(ActivatedRoute);
  private orderService = inject(OrderService);

  order = signal<Order | null>(null);

  private poolingInterval: any;

  // Status steps for the UI
  steps = ['PENDING', 'PREPARING', 'READY_FOR_PICKUP', 'IN_DELIVERY', 'DELIVERED'];

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.loadData(id);

      // Set up polling every 10 seconds
      this.poolingInterval = setInterval(() => {
        this.loadData(id);
      }, 5000);
    }
  }

  ngOnDestroy() {
    if (this.poolingInterval) {
      clearInterval(this.poolingInterval);
    }
  }
  loadData(id: string) {
    this.orderService.getMyOrders().subscribe(res => {
      console.log('Polled order data:', res);
      const found = res.orders?.find(o => o.id === id);
      if (found) this.order.set(found);
      if (found?.status === 'DELIVERED' || found?.status === 'CANCELLED') {
        clearInterval(this.poolingInterval);
      }
    });
  }

  // 1. Calculate when food is ready (Created Time + Prep Minutes)
  getReadyTime(order: Order): Date {
    if (!order.createdAt) return new Date();

    const startTime = new Date(order.createdAt);
    const prepMinutes = order.preparingFoodDeliveryTime || 0;

    // Add minutes to start time
    return new Date(startTime.getTime() + (prepMinutes * 60000));
  }

  // 2. Calculate Final Arrival (Ready Time + Driving Minutes)
  getFinalETA(order: Order): Date {
    const readyTime = this.getReadyTime(order);
    const driveMinutes = order.deliveryDurationTime || 0;

    // Add driving time to ready time
    return new Date(readyTime.getTime() + (driveMinutes * 60000));
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