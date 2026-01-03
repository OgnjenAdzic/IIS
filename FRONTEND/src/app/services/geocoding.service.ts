import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { map, of, tap } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class GeocodingService {
    private http = inject(HttpClient);
    private readonly API_URL = 'https://nominatim.openstreetmap.org/search';

    // Simple in-memory cache to make repeated searches instant
    private cache = new Map<string, any[]>();

    // NOVI SAD BOUNDING BOX (approximate coordinates)
    // [minLon, minLat, maxLon, maxLat]
    private readonly NOVI_SAD_VIEWBOX = '19.7,45.2,20.0,45.35';

    searchAddress(query: string) {
        // 1. Check Cache first
        if (this.cache.has(query)) {
            return of(this.cache.get(query)!);
        }

        return this.http.get<any[]>(this.API_URL, {
            params: {
                q: query,
                format: 'json',
                limit: '5',
                addressdetails: '1',
                'accept-language': 'sr-Latn, en', // Latin script
                // --- OPTIMIZATION PARAMETERS ---
                countrycodes: 'rs',      // Only Serbia
                viewbox: this.NOVI_SAD_VIEWBOX, // Only Novi Sad area
                bounded: '1'             // Strictly enforce the viewbox
            }
        }).pipe(
            map(results => results.map(item => {
                const addr = item.address || {};
                // Clean address logic
                const street = addr.road || addr.pedestrian || addr.street || '';
                const number = addr.house_number || '';
                // Only show city part if it's not implied
                let cleanAddress = `${street} ${number}`.trim();

                return {
                    display_name: cleanAddress + ', Novi Sad', // Enforce NS branding
                    lat: parseFloat(item.lat),
                    lon: parseFloat(item.lon)
                };
            })),
            // 2. Save to Cache
            tap(results => this.cache.set(query, results))
        );
    }
}