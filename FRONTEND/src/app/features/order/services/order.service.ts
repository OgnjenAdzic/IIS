import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';

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
}