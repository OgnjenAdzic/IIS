export interface FeeConfig {
    id?: string;
    itemRevenuePercent: number;
    deliveryRevenuePercent: number;
    appFeePercent: number;
    appFeeCap: number;
    smallOrderThreshold: number;
    smallOrderFee: number;
}

export interface ProfitLogItem {
    orderId: string;
    restaurantId: string;
    appFee: number;
    smallOrderFee: number;
    profitFromItems: number;
    profitFromDelivery: number;
    totalProfit: number;
    createdAt: string;
}
