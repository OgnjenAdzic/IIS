export interface FeeConfig {
    id?: string;
    itemRevenuePercent: number;
    deliveryRevenuePercent: number;
    appFeePercent: number;
    appFeeCap: number;
    smallOrderThreshold: number;
    smallOrderFee: number;
}