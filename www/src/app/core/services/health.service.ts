import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { HealthOverviewOutput, HealthRoutesOutput } from '../../features/dashboard/dashboard.models';

@Injectable({ providedIn: 'root' })
export class HealthService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getOverview(windowMinutes: number): Observable<HealthOverviewOutput> {
    const params = new HttpParams().set('window_minutes', String(windowMinutes));
    return this.http.get<HealthOverviewOutput>(`${this.apiUrl}/health`, { params });
  }

  getRoutes(windowMinutes: number, limit: number): Observable<HealthRoutesOutput> {
    const params = new HttpParams()
      .set('window_minutes', String(windowMinutes))
      .set('limit', String(limit));
    return this.http.get<HealthRoutesOutput>(`${this.apiUrl}/health/routes`, { params });
  }
}
