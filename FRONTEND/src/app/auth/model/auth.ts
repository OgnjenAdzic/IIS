export enum UserRole {
    ADMIN = 'ADMIN',
    CUSTOMER = 'CUSTOMER',
    DELIVERY_PERSON = 'DELIVERY_PERSON',
    RESTAURANT_WORKER = 'RESTAURANT_WORKER'
}

export interface User {
    id: string;
    username: string;
    role: UserRole;
    exp?: number;
}

export interface AuthResponse {
    token: string;
}

export interface RegisterRequest {
    username: string;
    password: string;
    role: UserRole;
}

export interface LoginRequest {
    username: string;
    password: string;
}