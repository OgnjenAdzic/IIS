export interface OrderItem {
    name: string;
    price: number;
    quantity: number;
}

export interface Order {
    id: string;
    restaurantId: string;
    customerId: string;
    status: string; // 'PENDING', 'PREPARING', etc.
    deliveryAddress: string;

    // Financials
    itemsTotal: number;
    deliveryFee: number;
    appFee: number;
    smallOrderFee: number;
    totalPrice: number;

    items: OrderItem[];
    createdAt: string;
    isPriority: boolean;

    estimatedDeliveryTime?: string; // New field for estimated delivery time
}