import { Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login';
import { RegisterComponent } from './auth/register/register';
import { Admin } from './landing/admin/admin';
import { Customer } from './landing/customer/customer';
import { DeliveryPerson } from './landing/delivery-person/delivery-person';
import { roleGuard } from './core/role.guard';
import { UserRole } from './auth/model/auth';
import { Restaurant } from './features/restaurant/pages/create-restaurant/restaurant';
import { RestaurantWorker } from './landing/restaurant-worker/restaurant-worker';
import { ManageRestaurant } from './features/restaurant/pages/manage-restaurant/manage-restaurant';

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
    {
        path: 'create-restaurant',
        component: Restaurant,
        canActivate: [roleGuard],
        data: { role: UserRole.ADMIN }
    },
    {
        path: 'manage-restaurant/:id',
        component: ManageRestaurant,
        canActivate: [roleGuard],
        data: { role: UserRole.RESTAURANT_WORKER }
    }

];
