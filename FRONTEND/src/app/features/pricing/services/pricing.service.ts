import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import { ConfigResponse, PricingRules, SystemStatus } from '../models/pricing.model';

@Injectable({
  providedIn: 'root'
})
export class PricingService {

  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/pricing`;

  calculatePrice(arg0: { customerLat: number; customerLon: number; restaurantLat: any; restaurantLon: any; isPriority: boolean; }) {
    return this.http.post<any>(`${this.apiUrl}/calculate`, arg0);
  }

  getConfig() {
    return this.http.get<ConfigResponse>(`${this.apiUrl}/config`);
  }

  updateStatus(isRushHour: boolean, isBadWeather: boolean) {
    return this.http.put<SystemStatus>(`${this.apiUrl}/status`, { isRushHour, isBadWeather });
  }

  updateRules(rules: PricingRules) {
    return this.http.post<PricingRules>(`${this.apiUrl}/rules`, rules);
  }
}