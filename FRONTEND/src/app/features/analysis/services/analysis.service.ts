import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/enviorment';
import { FeeConfig } from '../models/analysis.model';

@Injectable({
  providedIn: 'root'
})
export class AnalysisService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/analysis`;

  analysisUpdated = signal<number>(0);

  notifyChange() {
    console.log("Broadcasting Analysis Change...");
    this.analysisUpdated.update(v => v + 1);
  }

  getConfig() {
    return this.http.get<FeeConfig>(`${this.apiUrl}/config`);
  }

  updateConfig(config: FeeConfig) {
    return this.http.post<FeeConfig>(`${this.apiUrl}/config`, config);
  }

  calculateFees(itemsTotal: number, deliveryPrice: number) {
    return this.http.post<any>(`${this.apiUrl}/calculate`, { itemsTotal, deliveryPrice });
  }
}