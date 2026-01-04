// 1. The Item inside the cart (Matches "message CartItem" in Proto)
export interface CartItem {
    id: string;          // Unique ID of this cart entry
    menuItemId: string;  // Reference to the food item
    name: string;
    price: number;
    quantity: number;
}

// 2. The Full Response from Backend (Matches "message CartResponse")
export interface CartResponse {
    id: string;
    userId: string;
    restaurantId: string;
    items: CartItem[];
    totalPrice: number;
}

// 3. Payload for Adding Items (Matches "message AddItemRequest")
export interface AddItemRequest {
    restaurantId: string;
    menuItemId: string;
    name: string;
    price: number;
    quantity: number;
}

// 4. Payload for Updating (Matches "message UpdateQuantityRequest")
export interface UpdateQuantityRequest {
    itemId: string;
    quantity: number;
}