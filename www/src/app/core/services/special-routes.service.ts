import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SpecialRouteListOutput, SpecialRouteOutput, SpecialRouteType } from '../../features/dashboard/dashboard.models';

export interface SpecialRouteInput {
  path: string;
  max_distinct_requests: number;
  window_seconds: number;
  ban_seconds: number;
  enabled: boolean;
}

@Injectable({ providedIn: 'root' })
export class SpecialRoutesService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(type: SpecialRouteType): Observable<SpecialRouteListOutput> {
    return this.http.get<SpecialRouteListOutput>(`${this.apiUrl}/special-routes/${type}`);
  }

  create(type: SpecialRouteType, data: SpecialRouteInput): Observable<SpecialRouteOutput> {
    return this.http.post<SpecialRouteOutput>(`${this.apiUrl}/special-routes/${type}`, data);
  }

  update(type: SpecialRouteType, id: number, data: SpecialRouteInput): Observable<SpecialRouteOutput> {
    return this.http.put<SpecialRouteOutput>(`${this.apiUrl}/special-routes/${type}/${id}`, data);
  }

  delete(type: SpecialRouteType, id: number): Observable<SpecialRouteOutput> {
    return this.http.delete<SpecialRouteOutput>(`${this.apiUrl}/special-routes/${type}/${id}`);
  }
}
