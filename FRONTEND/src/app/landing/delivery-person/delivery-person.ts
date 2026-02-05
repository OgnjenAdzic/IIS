import { Component, inject, OnInit, signal, effect, OnDestroy, computed } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { OrderService } from '../../features/order/services/order.service';
import { Stakeholders } from '../../features/stakeholders/service/stakeholders';
import { Order } from '../../features/order/models/order.model';
import { AnalysisService } from '../../features/analysis/services/analysis.service';
import { FeeConfig } from '../../features/analysis/models/analysis.model';

@Component({
  selector: 'app-delivery-person',
  standalone: true,
  imports: [CommonModule, DatePipe, FormsModule],
  templateUrl: './delivery-person.html',
  styleUrl: './delivery-person.css'
})
export class DeliveryPerson implements OnInit, OnDestroy {
  orderService = inject(OrderService);
  stakeholdersService = inject(Stakeholders);
  analysisService = inject(AnalysisService)

  // Signals
  availableOrders = signal<Order[]>([]);
  currentJob = signal<Order | null>(null);
  analysisConfig = signal<FeeConfig | null>(null);

  bidInputs: { [id: string]: number } = {};
  myVehicle = '';
  private refreshInterval: any;
  private refreshDeliveryFee: any;



  constructor() {
    effect(() => {
      const profile = this.stakeholdersService.currentProfile();
      if (profile && profile.vehicle) {
        this.myVehicle = profile.vehicle; // e.g. "CAR", "BIKE"
        console.log("Vehicle identified:", this.myVehicle);
      }
    });
  }

  visibleOrders = computed(() => {
    const orders = this.availableOrders();
    const profile = this.stakeholdersService.currentProfile();
    const myVehicle = profile.vehicle;
    console.log("Filtered visible orders:", orders);
    return orders.filter(order => {
      if (!order.isPriority) return true;

      return myVehicle === 'CAR';
    });

  });

  ngOnInit() {
    this.stakeholdersService.loadCurrentProfile();

    this.getDeliveryRevenue();
    this.refreshData();

    this.refreshInterval = setInterval(() => {
      this.refreshData();
    }, 45000);

    this.refreshDeliveryFee = setInterval(() => {
      this.getDeliveryRevenue();
    }, 5000);
  }

  ngOnDestroy() {
    if (this.refreshInterval) {
      clearInterval(this.refreshInterval);
    }
    if (this.refreshDeliveryFee) {
      clearInterval(this.refreshDeliveryFee);
    }
  }

  getDeliveryRevenue() {
    this.analysisService.getConfig().subscribe({
      next: (res) => {
        this.analysisConfig.set(res);
        console.log("Company takes:", res.deliveryRevenuePercent + "%");
      },
      error: (err) => console.log("Failed to load fees: ", err)
    });
  }

  refreshData() {
    // 1. Load Available Orders (Bidding list)
    this.orderService.getAvailableOrders().subscribe({
      next: (res) => {
        this.availableOrders.set(res.orders || [])
      }
    });

    // 2. Load My Active Job (The new endpoint)
    this.orderService.getActiveJob().subscribe({
      next: (order) => {
        // Success: We have a job
        this.currentJob.set(order);
        this.availableOrders.set([]); // Clear bidding list
        console.log("Active job loaded:", order);
      },
      error: (err) => {
        // 404 Not Found means the driver is free
        // Set signal to null so the UI shows the Bidding list
        this.currentJob.set(null);
        this.orderService.getAvailableOrders().subscribe({
          next: (res) => {
            this.availableOrders.set(res.orders || [])
          }
        });
      }
    });
  }

  placeBid(order: Order) {
    const mins = this.bidInputs[order.id];
    if (!mins) return;

    this.orderService.bidForOrder(order.id, mins).subscribe({
      next: () => {
        alert("Bid Sent! If you win, it will appear in 'Current Job' shortly.");
        this.availableOrders.update(list => list.filter(o => o.id !== order.id));
      },
      error: (err) => alert("Failed to bid: " + err.error?.message || "Unknown error")
    });
  }

  // Actions for Active Job
  confirmPickup(id: string) {
    this.orderService.updateStatus(id, 'IN_DELIVERY').subscribe(() => this.refreshData());
  }

  completeDelivery(id: string) {
    this.orderService.updateStatus(id, 'DELIVERED').subscribe(() => this.refreshData());
  }

  getReadyTime(order: any): Date {
    if (!order.createdAt) return new Date();

    const created = new Date(order.createdAt);
    // Use the field name exactly as it appears in your console log
    const minutesToAdd = order.preparingFoodDeliveryTime || 0;

    return new Date(created.getTime() + (minutesToAdd * 60 * 1000));
  }

  calculateEarnings(deliveryFee: number): number {
    const config = this.analysisConfig();
    if (!config) {
      return 65
    }
    const driverSharePercent = 100 - config.deliveryRevenuePercent;
    return Math.round(deliveryFee * (driverSharePercent / 100));

  }
}