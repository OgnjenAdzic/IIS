import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login';
import { RegisterComponent } from './auth/register/register';
import { Admin } from './landing/admin/admin';
import { Customer } from './landing/customer/customer';
import { DeliveryPerson } from './landing/delivery-person/delivery-person';
import { RestaurantWorker } from './landing/restaurant-worker/restaurant-worker';
import { roleGuard } from './core/role.guard';
import { UserRole } from './auth/model/auth';

export const routes: Routes = [
    { path: '', redirectTo: 'login', pathMatch: 'full' },
    { path: 'login', component: LoginComponent },
    { path: 'register', component: RegisterComponent },
    {
        path: 'customer',
        component: Customer,
        canActivate: [roleGuard],
        data: { role: UserRole.CUSTOMER }
    }, {
        path: 'admin',
        component: Admin,
        canActivate: [roleGuard],
        data: { role: UserRole.ADMIN }
    },
    {
        path: 'delivery',
        component: DeliveryPerson,
        canActivate: [roleGuard],
        data: { role: UserRole.DELIVERY_PERSON }
    },
    {
        path: 'restaurant',
        component: RestaurantWorker,
        canActivate: [roleGuard],
        data: { role: UserRole.RESTAURANT_WORKER }
    },

];
