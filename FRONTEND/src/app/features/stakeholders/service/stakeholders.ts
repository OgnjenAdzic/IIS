import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import { AuthService } from '../../../auth/service/auth';
import { UserRole } from '../../../auth/model/auth';
import { catchError, map, of, Observable } from 'rxjs';

import {
  CreateCustomerRequest,
  CreateDeliveryPersonRequest,
  CustomerProfile,
  DeliveryPersonProfile,
  CreateRestaurantWorkerRequest,
  RestaurantWorkerPersonResponse,
  AllRestaurantWorkersResponse
} from '../models/stakeholder';

@Injectable({
  providedIn: 'root',
})
export class Stakeholders {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private apiUrl = `${environment.apiUrl}/stakeholders`;

  currentProfile = signal<any>(null);

  loadCurrentProfile() {
    const user = this.authService.currentUser();
    if (!user) return;

    // Determine endpoint based on role
    let endpoint = 'customer';
    if (user.role === UserRole.DELIVERY_PERSON) endpoint = 'delivery-person';
    if (user.role === UserRole.RESTAURANT_WORKER) endpoint = 'worker';

    // Fetch and SET the signal
    this.http.get(`${this.apiUrl}/${endpoint}/${user.id}`).subscribe({
      next: (profile) => {
        console.log("Caching Profile:", profile);
        this.currentProfile.set(profile); // <--- THIS STORES IT
      },
      error: (err) => console.error("Failed to load profile", err)
    });
  }

  hasProfile(): Observable<boolean> {
    const user = this.authService.currentUser();
    if (!user) return of(false);

    let endpoint = 'customer';
    if (user.role === UserRole.DELIVERY_PERSON) endpoint = 'delivery-person';
    if (user.role === UserRole.RESTAURANT_WORKER) endpoint = 'worker';

    return this.http.get(`${this.apiUrl}/${endpoint}/${user.id}`).pipe(
      map(() => true),
      catchError(() => of(false))
    );
  }

  createCustomerProfile(data: CreateCustomerRequest) {
    const user = this.authService.currentUser();
    const payload: CreateCustomerRequest = { ...data, userId: user?.id };
    return this.http.post<CustomerProfile>(`${this.apiUrl}/customer`, payload);
  }

  createDeliveryProfile(data: CreateDeliveryPersonRequest) {
    const user = this.authService.currentUser();
    const payload: CreateDeliveryPersonRequest = { ...data, userId: user?.id };
    return this.http.post<DeliveryPersonProfile>(`${this.apiUrl}/delivery-person`, payload);
  }

  createWorkerProfile(data: CreateRestaurantWorkerRequest) {
    const user = this.authService.currentUser();
    const payload: CreateRestaurantWorkerRequest = { ...data, userId: user?.id };
    return this.http.post<RestaurantWorkerPersonResponse>(`${this.apiUrl}/worker`, payload);
  }

  getAllWorkers() {
    return this.http.get<AllRestaurantWorkersResponse>(`${this.apiUrl}/workers`);
  }

  getCustomerProfile(id: string) {
    return this.http.get<CustomerProfile>(`${this.apiUrl}/customer/${id}`);
  }

  getCoordinates() {
    const profile = this.currentProfile();
    if (profile && profile.latitude && profile.longitude) {
      return { lat: profile.latitude, lon: profile.longitude };
    }
    return null;
  }
}
