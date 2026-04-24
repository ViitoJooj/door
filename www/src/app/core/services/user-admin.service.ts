import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AdminUser, UserListOutput, AdminCreateUserOutput, UserOutput } from '../../features/dashboard/dashboard.models';

export interface CreateUserInput {
  username: string;
  email: string;
  role: string;
}

export interface UpdateUserInput {
  username: string;
  email: string;
  password: string;
  role: string;
  active: boolean;
}

@Injectable({ providedIn: 'root' })
export class UserAdminService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(): Observable<UserListOutput> {
    return this.http.get<UserListOutput>(`${this.apiUrl}/users`);
  }

  create(data: CreateUserInput): Observable<AdminCreateUserOutput> {
    return this.http.post<AdminCreateUserOutput>(`${this.apiUrl}/users`, data);
  }

  update(id: number, data: UpdateUserInput): Observable<UserOutput> {
    return this.http.put<UserOutput>(`${this.apiUrl}/users/${id}`, data);
  }

  delete(id: number): Observable<{ success: boolean; message: string }> {
    return this.http.delete<{ success: boolean; message: string }>(`${this.apiUrl}/users/${id}`);
  }
}
