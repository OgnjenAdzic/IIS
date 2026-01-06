import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import { Order } from '../models/order.model';

@Injectable({
  providedIn: 'root'
})
export class OrderService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/orders`;

  createOrder(data: {
    userId?: string,
    customAddress?: string,
    customLat?: number,
    customLon?: number,
    isPriority: boolean
  }) {
    return this.http.post(this.apiUrl, data);
  }

  getMyOrders() {
    return this.http.get<{ orders: Order[] }>(this.apiUrl);
  }
  getRestaurantOrders(restaurantId: string, status: string = '') {
    let url = `${this.apiUrl}/restaurant/${restaurantId}`;

    if (status) {
      url += `?status=${status}`;
    }

    return this.http.get<{ orders: Order[] }>(url);
  }

  updateStatus(orderId: string, status: string, estimatedDeliveryMinutes: number = 0) {
    return this.http.put(`${this.apiUrl}/${orderId}/status`, {
      orderId,
      status,
      estimatedDeliveryMinutes
    });
  }
}

