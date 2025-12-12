export interface CreateRestaurantRequest {
    name: string;
    category: string;
    address: string;
    latitude: number;
    longitude: number;
}

export interface MenuItem {
    id: string;
    name: string;
    price: number;
}

export interface Menu {
    id: string;
    items: MenuItem[];
}

export interface Restaurant {
    id: string;
    name: string;
    category: string;
    isOpen: boolean;
    menu: Menu;
    address: string;
    latitude: number;
    longitude: number;
}

export interface GetAllRestaurantsResponse {
    restaurants: Restaurant[];
}