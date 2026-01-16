import { Component, inject, ChangeDetectorRef, OnInit } from '@angular/core';
import { AuthService } from '../../auth/service/auth';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { PricingService } from '../../features/pricing/services/pricing.service';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { SystemStatus } from '../../features/pricing/models/pricing.model';
import { AnalysisService } from '../../features/analysis/services/analysis.service';

@Component({
  selector: 'app-admin',
  imports: [RouterLink, ReactiveFormsModule, CommonModule],
  templateUrl: './admin.html',
  styleUrl: './admin.css',
})
export class Admin implements OnInit {
  authService = inject(AuthService);
  pricingService = inject(PricingService);
  analysisService = inject(AnalysisService);
  fb = inject(FormBuilder);
  private cdr = inject(ChangeDetectorRef);

  status: SystemStatus = { isRushHour: false, isBadWeather: false };
  isLoading = true;
  successMessage = '';

  pricingForm = this.fb.group({
    basePrice: [0, [Validators.required, Validators.min(0)]],
    pricePerKm: [0, [Validators.required, Validators.min(0)]],
    rushHourFee: [0, [Validators.required, Validators.min(0)]],
    weatherFee: [0, [Validators.required, Validators.min(0)]]
  });

  analysisForm = this.fb.group({
    itemRevenuePercent: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
    deliveryRevenuePercent: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
    appFeePercent: [0, [Validators.required, Validators.min(0), Validators.max(100)]],
    appFeeCap: [0, [Validators.required, Validators.min(0)]],
    smallOrderThreshold: [0, [Validators.required, Validators.min(0)]],
    smallOrderFee: [0, [Validators.required, Validators.min(0)]]
  });

  ngOnInit() {
    this.loadConfig();
    this.loadAnalysisConfig();
  }

  loadConfig() {
    this.pricingService.getConfig().subscribe({
      next: (res) => {
        this.status = res.status;
        console.log(res.status);
        this.pricingForm.patchValue(res.rules);
        this.isLoading = false;
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load config', err);
        this.isLoading = false;
      }
    });
  }

  loadAnalysisConfig() {
    this.analysisService.getConfig().subscribe({
      next: (config) => {
        this.analysisForm.patchValue(config);
        console.log(config);
        this.cdr.detectChanges();
      },
      error: (err) => {
        console.error('Failed to load analysis config', err);
      }
    });
  }

  // --- ACTIONS ---

  toggleRushHour() {
    const newStatus = !this.status.isRushHour;
    // Optimistic Update (Update UI immediately)
    this.status.isRushHour = newStatus;

    this.pricingService.updateStatus(newStatus, this.status.isBadWeather).subscribe({
      next: () => {
        this.cdr.detectChanges();
        this.pricingService.notifyChange();
      },
      error: () => this.status.isRushHour = !newStatus // Revert on error
    });
  }

  toggleWeather() {
    const newStatus = !this.status.isBadWeather;
    this.status.isBadWeather = newStatus;

    this.pricingService.updateStatus(this.status.isRushHour, newStatus).subscribe({
      next: () => {
        this.cdr.detectChanges();
        this.pricingService.notifyChange();
      },
      error: () => this.status.isBadWeather = !newStatus
    });
  }

  saveRules() {
    if (this.pricingForm.valid) {
      const formVal = this.pricingForm.value as any;

      this.pricingService.updateRules(formVal).subscribe({
        next: () => {
          this.showSuccess('Pricing rules updated successfully!');
          this.cdr.detectChanges();
          this.pricingService.notifyChange();
        },
        error: (err) => alert('Failed to update rules')
      });
    }
  }

  saveAnalysisRules() {
    if (this.analysisForm.valid) {
      this.analysisService.updateConfig(this.analysisForm.value as any).subscribe({
        next: () => {
          this.showSuccess('Financial config updated!');
          this.cdr.detectChanges();
          this.pricingService.notifyChange();
        },
        error: () => alert('Failed to update analysis config')
      });
    }
  }

  private showSuccess(msg: string) {
    this.successMessage = msg;
    setTimeout(() => this.successMessage = '', 3000);
  }
}
