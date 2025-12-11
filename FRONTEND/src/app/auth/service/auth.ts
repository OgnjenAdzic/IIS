import { Injectable, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/enviorment';
import { AuthResponse, LoginRequest, RegisterRequest, User, UserRole } from '../model/auth';
import { tap } from 'rxjs';
import { jwtDecode } from 'jwt-decode';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = `${environment.apiUrl}/auth`;
  private tokenKey = 'delivery_token';

  // Signal to hold the current user state
  currentUser = signal<User | null>(this.getUserFromStorage());

  // Computed signal to check if logged in
  isLoggedIn = computed(() => !!this.currentUser());

  constructor(private http: HttpClient, private router: Router) { }

  login(credentials: LoginRequest) {
    return this.http.post<AuthResponse>(`${this.apiUrl}/login`, credentials).pipe(
      tap(response => {
        if (response.token) {
          localStorage.setItem(this.tokenKey, response.token);
          const user = this.decodeToken(response.token);
          this.currentUser.set(user);
        }
      })
    );
  }

  register(data: RegisterRequest) {
    return this.http.post(`${this.apiUrl}/register`, data);
  }

  logout() {
    localStorage.removeItem(this.tokenKey);
    this.currentUser.set(null);
    this.router.navigate(['/login']);
  }

  getToken(): string | null {
    return localStorage.getItem(this.tokenKey);
  }

  private getUserFromStorage(): User | null {
    const token = this.getToken();
    return token ? this.decodeToken(token) : null;
  }

  private decodeToken(token: string): User | null {
    try {
      const decoded: any = jwtDecode(token);
      return {
        id: decoded.sub,
        username: decoded.username,
        role: decoded.role as UserRole,
        exp: decoded.exp
      };
    } catch (e) {
      return null;
    }
  }
}