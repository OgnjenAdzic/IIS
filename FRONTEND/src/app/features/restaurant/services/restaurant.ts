import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { GeocodingService } from '../../../services/geocoding.service';
import { CreateRestaurantRequest } from '../models/restaurant';
import { environment } from '../../../../environments/enviorment';
import { GetAllRestaurantsResponse, Restaurant } from '../models/restaurant';

@Injectable({
  providedIn: 'root',
})
export class RestaurantService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/restaurant`;

  getAll() {
    return this.http.get<GetAllRestaurantsResponse>(this.apiUrl);
  }

  getById(id: string) {
    return this.http.get<Restaurant>(`${this.apiUrl}/${id}`);
  }

  createRestaurant(data: any) {
    return this.http.post(this.apiUrl, data);
  }

  updateStatus(id: string, isOpen: boolean) {
    return this.http.put(`${this.apiUrl}/${id}/status`, { id, isOpen });
  }

  addMenuItem(restaurantId: string, name: string, price: number) {
    return this.http.post(`${this.apiUrl}/${restaurantId}/menu`, { restaurantId, name, price });
  }

  deleteMenuItem(itemId: string) {
    return this.http.delete(`${this.apiUrl}/menu/${itemId}`);
  }

  updateItemPrice(itemId: string, price: number) {
    return this.http.put(`${this.apiUrl}/menu/${itemId}/price`, { id: itemId, price });
  }
}
