import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import {
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  RefreshRequest,
  RefreshResponse
} from '../../features/auth/auth.models';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = environment.apiUrl;
  // access token fica só em memória, não vai pro localStorage
  private accessToken: string | null = null;

  constructor(private http: HttpClient, private router: Router) {}

  login(data: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/auth/login`, data).pipe(
      tap(res => this.saveTokens(res))
    );
  }

  register(data: RegisterRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.apiUrl}/auth/register`, data).pipe(
      tap(res => this.saveTokens(res))
    );
  }

  refresh(): Observable<RefreshResponse> {
    const refreshToken = this.getRefreshToken();
    const body: RefreshRequest = { refresh_token: refreshToken! };
    return this.http.post<RefreshResponse>(`${this.apiUrl}/auth/token`, body).pipe(
      tap(res => this.setAccessToken(res.access_token))
    );
  }

  logout(): void {
    const refreshToken = this.getRefreshToken();
    if (refreshToken) {
      this.http.post(`${this.apiUrl}/auth/logout`, { refresh_token: refreshToken }).subscribe();
    }
    this.clearTokens();
    this.router.navigate(['/auth']);
  }

  getAccessToken(): string | null {
    return this.accessToken;
  }

  setAccessToken(token: string): void {
    this.accessToken = token;
  }

  getRefreshToken(): string | null {
    return localStorage.getItem('refresh_token');
  }

  isAuthenticated(): boolean {
    return !!this.accessToken;
  }

  private saveTokens(res: AuthResponse): void {
    this.accessToken = res.access_token;
    localStorage.setItem('refresh_token', res.refresh_token);
  }

  private clearTokens(): void {
    this.accessToken = null;
    localStorage.removeItem('refresh_token');
  }
}
