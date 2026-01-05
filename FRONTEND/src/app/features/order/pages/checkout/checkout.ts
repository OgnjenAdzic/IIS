import { Component, inject, OnInit, signal, effect } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';

// Services
import { CartService } from '../../../cart/services/services';
import { Stakeholders } from '../../../stakeholders/service/stakeholders';
import { PricingService } from '../../../pricing/services/pricing.service';
import { AnalysisService } from '../../../analysis/services/analysis.service';
import { OrderService } from '../../services/order.service';
import { GeocodingService } from '../../../../services/geocoding.service';
import { RestaurantService } from '../../../restaurant/services/restaurant';

@Component({
  selector: 'app-checkout',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './checkout.html',
  styleUrl: './checkout.css'
})
export class Checkout implements OnInit {
  // Injections
  cartService = inject(CartService);
  stakeholdersService = inject(Stakeholders);
  pricingService = inject(PricingService);
  analysisService = inject(AnalysisService);
  orderService = inject(OrderService);
  geoService = inject(GeocodingService);
  restaurantService = inject(RestaurantService);
  router = inject(Router);

  // State
  cartItems = this.cartService.cartItems;
  cartTotal = this.cartService.totalPrice; // Items Subtotal

  // Address State
  useCustomAddress = signal<boolean>(false);
  defaultProfile: any = null;
  customAddress = { name: '', lat: 0, lon: 0 };
  addressSuggestions: any[] = [];

  // Financial Breakdown State
  fees = signal({
    baseDelivery: 0,
    distanceFee: 0,
    rushHourFee: 0,
    weatherFee: 0,
    appFee: 0,
    smallOrderFee: 0,
    deliveryTotal: 0, // Base + Dist + Rush + Weather
    grandTotal: 0
  });

  isPriority = signal<boolean>(false);
  loadingCosts = signal<boolean>(false);
  restaurantLocation: { lat: number, lon: number } | null = null;

  ngOnInit() {
    // 1. Check if Cart is empty
    if (this.cartItems().length === 0) {
      this.router.navigate(['/customer']);
      return;
    }

    // 2. Load Customer Profile (Default Address)
    this.stakeholdersService.loadCurrentProfile(); // Refresh logic
    // We assume currentProfile is populated or we fetch it
    const restaurantId = this.cartService.cartRestaurantId();

    if (!restaurantId) {
      console.error("No Restaurant ID found in cart");
      return;
    }
    // 3. Load Restaurant Location (Need this for distance)
    this.restaurantService.getById(restaurantId).subscribe(res => {
      this.restaurantLocation = { lat: res.latitude, lon: res.longitude };

      const cached = this.stakeholdersService.currentProfile();
      if (cached) {
        this.defaultProfile = cached;
        this.recalculateCosts();
      }
    });
  }

  // --- ADDRESS SEARCH (Reuse Logic) ---
  onAddressInput(event: any) {
    const query = event.target.value;
    if (query.length > 2) {
      this.geoService.searchAddress(query).subscribe(res => this.addressSuggestions = res);
    }
  }

  selectCustomAddress(addr: any) {
    this.customAddress = { name: addr.display_name, lat: addr.lat, lon: addr.lon };
    this.addressSuggestions = [];
    this.recalculateCosts(); // Recalc price for new location
  }

  toggleAddressMode() {
    this.useCustomAddress.update(v => !v);
    this.recalculateCosts();
  }

  togglePriority() {
    this.isPriority.update(v => !v);
    this.recalculateCosts();
  }

  // --- THE CORE CALCULATION LOGIC ---
  recalculateCosts() {
    if (!this.restaurantLocation) return;

    // 1. Determine active coordinates
    let userLat, userLon;
    if (this.useCustomAddress()) {
      if (!this.customAddress.lat) return; // Haven't selected yet
      userLat = this.customAddress.lat;
      userLon = this.customAddress.lon;
    } else {
      if (!this.defaultProfile) return;
      userLat = this.defaultProfile.latitude;
      userLon = this.defaultProfile.longitude;
    }

    this.loadingCosts.set(true);

    // 2. Call PRICING Service
    this.pricingService.calculatePrice({
      customerLat: userLat,
      customerLon: userLon,
      restaurantLat: this.restaurantLocation.lat,
      restaurantLon: this.restaurantLocation.lon,
      isPriority: this.isPriority()
    }).subscribe(priceRes => {

      // 3. Call ANALYSIS Service
      // We pass the Cart Total + The Delivery Price we just got
      this.analysisService.calculateFees(this.cartTotal(), priceRes.finalPrice).subscribe(analysisRes => {

        // 4. Update UI State
        this.fees.set({
          baseDelivery: priceRes.basePrice,
          distanceFee: priceRes.distancePrice,
          rushHourFee: priceRes.rushHourFee,
          weatherFee: priceRes.weatherFee,
          deliveryTotal: priceRes.finalPrice,

          appFee: analysisRes.appFee,
          smallOrderFee: analysisRes.smallOrderFee,

          // GRAND TOTAL = Items + Delivery + App Fees
          grandTotal: this.cartTotal() + priceRes.finalPrice + analysisRes.appFee + analysisRes.smallOrderFee
        });

        this.loadingCosts.set(false);
      });
    });
  }

  placeOrder() {
    const payload = {
      isPriority: this.isPriority(),
      customAddress: this.useCustomAddress() ? this.customAddress.name : '',
      customLat: this.useCustomAddress() ? this.customAddress.lat : 0,
      customLon: this.useCustomAddress() ? this.customAddress.lon : 0,
    };

    this.orderService.createOrder(payload).subscribe({
      next: (order) => {
        alert('Order Placed Successfully!');
        this.cartService.clearCart(); // Frontend cleanup
        this.router.navigate(['/customer']); // Or order history
      },
      error: (err) => alert('Failed to place order')
    });
  }
}