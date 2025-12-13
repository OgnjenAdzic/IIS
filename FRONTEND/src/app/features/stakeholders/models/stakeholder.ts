export type VehicleType = 'CAR' | 'BIKE' | 'SCOOTER';

export interface CreateCustomerRequest {
    userId?: string; // Optional because we attach it in the service
    firstName: string;
    lastName: string;
    address: string;
    latitude: number;
    longitude: number;
}

export interface CreateDeliveryPersonRequest {
    userId?: string;
    firstName: string;
    lastName: string;
    vehicle: VehicleType;
}

export interface CustomerProfile {
    userId: string;
    firstName: string;
    lastName: string;
    address: string;
    latitude: number;
    longitude: number;
}


export interface DeliveryPersonProfile {
    userId: string;
    firstName: string;
    lastName: string;
    vehicle: VehicleType;
    isWorking: boolean;
    deliveryCount: number;
}