import { Component, inject } from '@angular/core';
import { RestaurantService } from '../../services/restaurant';
import { GeocodingService } from '../../services/geocoding.service';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';

@Component({
  selector: 'app-restaurant',
  imports: [ReactiveFormsModule, CommonModule],
  templateUrl: './restaurant.html',
  styleUrl: './restaurant.css',
})
export class Restaurant {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private geoService = inject(GeocodingService);
  private restaurantService = inject(RestaurantService)

  addressSuggestions: any[] = [];

  restaurantForm = this.fb.group({
    name: ['', Validators.required],
    category: ['', Validators.required],
    address: ['', Validators.required],
    latitude: [0.0, Validators.required],
    longitude: [0.0, Validators.required]
  });

  onAddressInput(event: any) {
    const query = event.target.value;
    if (query.length > 3) {
      this.geoService.searchAddress(query).subscribe(results => {
        this.addressSuggestions = results;
      });
    } else {
      this.addressSuggestions = [];
    }
  }

  selectAddress(item: any) {
    console.log("Selected Item:", item);

    this.restaurantForm.patchValue({
      address: item.display_name,
      latitude: item.lat,
      longitude: item.lon
    });
    this.addressSuggestions = [];
  }

  onSubmit() {
    if (this.restaurantForm.valid) {
      const req = this.restaurantForm.value as any;

      this.restaurantService.createRestaurant(req).subscribe({
        next: () => {
          alert('Restaurant Created Successfully!');
          this.router.navigate(['/admin']);
        },
        error: (err) => {
          console.error(err);
          alert('Failed to create restaurant.');
        }
      });
    } else {
      alert("Please fill all fields and select a valid address.");
    }
  }
}
