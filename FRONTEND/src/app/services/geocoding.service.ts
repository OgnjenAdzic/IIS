import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class GeocodingService {
    private http = inject(HttpClient);
    private readonly API_URL = 'https://nominatim.openstreetmap.org/search';

    searchAddress(query: string) {
        return this.http.get<any[]>(this.API_URL, {
            params: {
                q: query,
                format: 'json',
                limit: '5',
                addressdetails: '1',     // <--- Request detailed parts (Street, Number, City)
                'accept-language': 'sr-Latn, en' // <--- Request Latin Script (Serbian Latin or English)
            }
        }).pipe(
            map(results => results.map(item => {
                // 1. Construct a clean address manually
                const addr = item.address || {};
                const street = addr.road || addr.pedestrian || addr.street || '';
                const number = addr.house_number || '';
                const city = addr.city || addr.town || addr.village || '';

                // Format: "Jevrejska 1, Novi Sad"
                let cleanAddress = `${street} ${number}, ${city}`;

                // Fallback if data is missing, use the first part of display_name
                if (!street) {
                    cleanAddress = item.display_name.split(',').slice(0, 3).join(',');
                }

                return {
                    // We return our clean string instead of the raw display_name
                    display_name: cleanAddress,
                    // 2. Ensure we parse coordinates correctly
                    lat: parseFloat(item.lat),
                    lon: parseFloat(item.lon)
                };
            }))
        );
    }
}