export interface PricingRules {
    id?: string;
    basePrice: number;
    pricePerKm: number;
    rushHourFee: number;
    weatherFee: number;
}

export interface SystemStatus {
    isRushHour: boolean;
    isBadWeather: boolean;
}

export interface ConfigResponse {
    rules: PricingRules;
    status: SystemStatus;
}