import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Application, ApplicationOutput } from '../../features/dashboard/dashboard.models';

@Injectable({
  providedIn: 'root'
})
export class ApplicationService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(): Observable<Application[]> {
    return this.http.get<Application[]>(`${this.apiUrl}/applications`);
  }

  create(data: { url: string; country: string }): Observable<ApplicationOutput> {
    return this.http.post<ApplicationOutput>(`${this.apiUrl}/applications`, data);
  }

  delete(id: number): Observable<ApplicationOutput> {
    return this.http.delete<ApplicationOutput>(`${this.apiUrl}/applications/${id}`);
  }
}
