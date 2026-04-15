import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { RequestLogListResponse } from '../../features/dashboard/dashboard.models';

@Injectable({
  providedIn: 'root'
})
export class RequestLogService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(): Observable<RequestLogListResponse> {
    return this.http.get<RequestLogListResponse>(`${this.apiUrl}/logs`);
  }
}
