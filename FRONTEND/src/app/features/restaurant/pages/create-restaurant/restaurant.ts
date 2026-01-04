import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { RestaurantService } from '../../services/restaurant';
import { GeocodingService } from '../../../../services/geocoding.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Subject, debounceTime, distinctUntilChanged, switchMap } from 'rxjs';
import { Stakeholders } from '../../../stakeholders/service/stakeholders';

@Component({
  selector: 'app-restaurant',
  standalone: true,
  imports: [ReactiveFormsModule, CommonModule],
  templateUrl: './restaurant.html',
  styleUrl: './restaurant.css',
})
export class Restaurant implements OnInit { // Implement OnInit
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private geoService = inject(GeocodingService);
  private restaurantService = inject(RestaurantService);
  private stakeholdersService = inject(Stakeholders);
  private cdr = inject(ChangeDetectorRef);

  addressSuggestions: any[] = [];
  availableManagers: any[] = [];

  // 1. Create a Subject to handle address input
  private addressSearchSubject = new Subject<string>();

  restaurantForm = this.fb.group({
    name: ['', Validators.required],
    category: ['', Validators.required],
    address: ['', Validators.required],
    latitude: [{ value: 0.0, disabled: true }, Validators.required], // Disable lat/lon initially
    longitude: [{ value: 0.0, disabled: true }, Validators.required], // Disable lat/lon initially
    managerId: ['', Validators.required]
  });

  ngOnInit(): void {
    this.stakeholdersService.getAllWorkers().subscribe({
      next: (res: any) => {
        console.log("Loaded workers:", res);
        this.availableManagers = res.workers || [];
        this.cdr.detectChanges();
      },
      error: (err) => console.error("Could not load workers", err)
    });
    // 2. Setup the reactive search pipeline
    this.addressSearchSubject.pipe(
      debounceTime(100),
      distinctUntilChanged(),
      switchMap(query => { // Cancel previous search if a new one comes in
        if (query.length > 1) { // Only search if query is long enough
          return this.geoService.searchAddress(query);
        } else {
          return []; // Return empty array if query is too short
        }
      })
    ).subscribe(results => {
      this.addressSuggestions = results;
    });
  }

  // 3. Update onAddressInput to use the Subject
  onAddressInput(event: Event) {
    const inputElement = event.target as HTMLInputElement;
    this.addressSearchSubject.next(inputElement.value);
  }

  selectAddress(item: any) {
    console.log("Selected Item:", item);

    this.restaurantForm.patchValue({
      address: item.display_name,
      latitude: item.lat,
      longitude: item.lon
    });

    // Enable lat/lon fields after selection
    this.restaurantForm.get('latitude')?.enable();
    this.restaurantForm.get('longitude')?.enable();

    this.addressSuggestions = []; // Clear suggestions
  }

  onSubmit() {
    if (this.restaurantForm.valid && this.restaurantForm.get('latitude')?.enabled) {
      const req = this.restaurantForm.getRawValue(); // <--- Use getRawValue()

      this.restaurantService.createRestaurant(req).subscribe({
        next: () => {
          alert('Restaurant Created Successfully!');
          this.router.navigate(['/admin']);
        },
        error: (err) => {
          console.error(err);
          // Better error message from backend if possible
          alert(`Failed to create restaurant: ${err.error?.message || 'Unknown error'}`);
        }
      });
    } else {
      alert("Please fill all required fields and select a valid address from the suggestions.");
    }
  }
}