import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from '../auth/service/auth';
import { UserRole } from '../auth/model/auth';

export const roleGuard: CanActivateFn = (route, state) => {
    const authService = inject(AuthService);
    const router = inject(Router);

    const user = authService.currentUser();
    const requiredRole = route.data['role'] as UserRole;
    const requiredRoles = route.data['roles'] as UserRole[];

    // 1. Check if logged in
    if (!user) {
        return router.createUrlTree(['/login']);
    }

    if (requiredRoles && requiredRoles.includes(user.role)) {
        return true;
    }

    // 4. Check Single (Allow if user's role matches exactly)
    if (requiredRole && user.role === requiredRole) {
        return true;
    }

    // 5. Unauthorized
    // Optional: Redirect them to their actual dashboard instead of just blocking
    if (user.role === UserRole.ADMIN) return router.createUrlTree(['/admin']);
    if (user.role === UserRole.CUSTOMER) return router.createUrlTree(['/customer']);


    // 3. If role is wrong, redirect to their correct dashboard
    // (Optional: or just show an error page)
    alert('You are not authorized to view this page.');
    return false;
};