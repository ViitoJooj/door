import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LoginRequest, RegisterRequest, AuthResponse, RegisterResponse, UserData } from '../../features/auth/auth.models';

const USER_STORAGE_KEY = 'ward_user';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = environment.apiUrl;
  private currentUser: UserData | null = null;

  constructor(private http: HttpClient, private router: Router) {
    const stored = localStorage.getItem(USER_STORAGE_KEY);
    if (stored) {
      try {
        this.currentUser = JSON.parse(stored);
      } catch {
        localStorage.removeItem(USER_STORAGE_KEY);
      }
    }
  }

  login(data: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/auth/login`, data).pipe(
      tap((res) => {
        if (res.data) {
          this.currentUser = res.data;
          localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(res.data));
        }
      })
    );
  }

  register(data: RegisterRequest): Observable<RegisterResponse> {
    return this.http.post<RegisterResponse>(`${this.apiUrl}/auth/register`, data).pipe(
      tap((res) => {
        if (res.data) {
          this.currentUser = res.data;
          localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(res.data));
        }
      })
    );
  }

  logout(): void {
    this.http.post(`${this.apiUrl}/auth/logout`, {}).subscribe();
    this.currentUser = null;
    localStorage.removeItem(USER_STORAGE_KEY);
    this.router.navigate(['/auth']);
  }

  isAuthenticated(): boolean {
    return this.currentUser !== null;
  }

  getCurrentUser(): UserData | null {
    return this.currentUser;
  }
}
