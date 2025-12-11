import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../auth/service/auth';
import { UserRole } from '../auth/model/auth';

export const roleGuard: CanActivateFn = (route, state) => {
    const authService = inject(AuthService);
    const router = inject(Router);

    const user = authService.currentUser();
    const expectedRole = route.data['role'] as UserRole;

    // 1. Check if logged in
    if (!user) {
        return router.createUrlTree(['/login']);
    }

    // 2. Check if role matches
    if (user.role === expectedRole) {
        return true;
    }

    // 3. If role is wrong, redirect to their correct dashboard
    // (Optional: or just show an error page)
    alert('You are not authorized to view this page.');
    return false;
};