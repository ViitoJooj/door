import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { RateLimitOutput } from '../../features/dashboard/dashboard.models';

export interface UpdateRateLimitInput {
  requests_per_second: number;
  burst: number;
  progressive_rate_limit: boolean;
}

@Injectable({ providedIn: 'root' })
export class RateLimitService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get(): Observable<RateLimitOutput> {
    return this.http.get<RateLimitOutput>(`${this.apiUrl}/rate-limit`);
  }

  update(data: UpdateRateLimitInput): Observable<RateLimitOutput> {
    return this.http.put<RateLimitOutput>(`${this.apiUrl}/rate-limit`, data);
  }
}
