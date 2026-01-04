import { Component, inject, ChangeDetectorRef } from '@angular/core';
import { AuthService } from '../../auth/service/auth';
import { RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { PricingService } from '../../features/pricing/services/pricing.service';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { SystemStatus } from '../../features/pricing/models/pricing.model';

@Component({
  selector: 'app-admin',
  imports: [RouterLink, ReactiveFormsModule, CommonModule],
  templateUrl: './admin.html',
  styleUrl: './admin.css',
})
export class Admin {
  authService = inject(AuthService);
  pricingService = inject(PricingService);
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

  ngOnInit() {
    this.loadConfig();
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

  // --- ACTIONS ---

  toggleRushHour() {
    const newStatus = !this.status.isRushHour;
    // Optimistic Update (Update UI immediately)
    this.status.isRushHour = newStatus;

    this.pricingService.updateStatus(newStatus, this.status.isBadWeather).subscribe({
      next: () => {
        this.cdr.detectChanges();
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
        },
        error: (err) => alert('Failed to update rules')
      });
    }
  }

  private showSuccess(msg: string) {
    this.successMessage = msg;
    setTimeout(() => this.successMessage = '', 3000); // Hide after 3s
  }
}
