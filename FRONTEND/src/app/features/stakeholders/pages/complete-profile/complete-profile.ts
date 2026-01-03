import { Component, inject, OnInit } from '@angular/core';
import { Subject, debounceTime, distinctUntilChanged, switchMap } from 'rxjs';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../../../auth/service/auth';
import { Stakeholders } from '../../service/stakeholders';
import { GeocodingService } from '../../../../services/geocoding.service';
import { UserRole } from '../../../../auth/model/auth';
import { CreateCustomerRequest, CreateDeliveryPersonRequest, VehicleType } from '../../models/stakeholder';

@Component({
  selector: 'app-complete-profile',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './complete-profile.html',
  styleUrl: './complete-profile.css',
})
export class CompleteProfile implements OnInit {
  fb = inject(FormBuilder);
  router = inject(Router);
  authService = inject(AuthService);
  stakeholdersService = inject(Stakeholders);
  geoService = inject(GeocodingService);

  userRole: UserRole | undefined;
  addressSuggestions: any[] = [];
  profileForm = this.fb.group({
    firstName: ['', Validators.required],
    lastName: ['', Validators.required],
    address: [''],
    latitude: [0],
    longitude: [0],
    vehicle: ['CAR']
  });

  private searchSubject = new Subject<string>();

  ngOnInit() {
    const user = this.authService.currentUser();
    this.userRole = user?.role;
    if (!user) this.router.navigate(['/login']);
    this.searchSubject.pipe(
      debounceTime(100),
      distinctUntilChanged(),
      switchMap(query => {
        // Cancel previous request if new one comes in
        return this.geoService.searchAddress(query);
      })
    ).subscribe(results => {
      this.addressSuggestions = results;
    });
  }

  onAddressInput(event: any) {
    const query = event.target.value;
    if (query.length > 2) {
      this.searchSubject.next(query);
    } else {
      this.addressSuggestions = [];
    }
  }

  selectAddress(item: any) {
    this.profileForm.patchValue({
      address: item.display_name,
      latitude: item.lat,
      longitude: item.lon
    });
    this.addressSuggestions = [];
  }

  onSubmit() {
    if (this.profileForm.invalid) return;

    const formVal = this.profileForm.getRawValue();

    if (this.userRole === UserRole.CUSTOMER) {
      const customerData: CreateCustomerRequest = {
        firstName: formVal.firstName!,
        lastName: formVal.lastName!,
        address: formVal.address!,
        latitude: formVal.latitude!,
        longitude: formVal.longitude!
      };

      this.stakeholdersService.createCustomerProfile(customerData).subscribe({
        next: () => this.router.navigate(['/customer']),
        error: (err) => console.error("Customer creation failed", err)
      });
    }
    else if (this.userRole === UserRole.DELIVERY_PERSON) {
      const deliveryData: CreateDeliveryPersonRequest = {
        firstName: formVal.firstName!,
        lastName: formVal.lastName!,
        vehicle: formVal.vehicle as VehicleType
      };

      this.stakeholdersService.createDeliveryProfile(deliveryData).subscribe({
        next: () => this.router.navigate(['/delivery']),
        error: (err) => console.error("Delivery profile failed", err)
      });
    }
    else if (this.userRole === UserRole.RESTAURANT_WORKER) {
      const workerData = {
        firstName: formVal.firstName!,
        lastName: formVal.lastName!
      };

      this.stakeholdersService.createWorkerProfile(workerData).subscribe({
        next: () => this.router.navigate(['/restaurant']),
        error: (err) => console.error("Worker profile creation failed", err)
      });
    }
  }
}
