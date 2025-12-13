import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import { AuthService } from '../../../auth/service/auth';
import { UserRole } from '../../../auth/model/auth';
import { catchError, map, of, Observable } from 'rxjs';

import {
  CreateCustomerRequest,
  CreateDeliveryPersonRequest,
  CustomerProfile,
  DeliveryPersonProfile
} from '../models/stakeholder';

@Injectable({
  providedIn: 'root',
})
export class Stakeholders {
  private http = inject(HttpClient);
  private authService = inject(AuthService);
  private apiUrl = `${environment.apiUrl}/stakeholders`;

  hasProfile(): Observable<boolean> {
    const user = this.authService.currentUser();
    if (!user) return of(false);

    const endpoint = user.role === UserRole.DELIVERY_PERSON
      ? 'delivery-person'
      : 'customer';

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
}
