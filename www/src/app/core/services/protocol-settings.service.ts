import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { ProtocolSettingsOutput } from '../../features/dashboard/dashboard.models';

export interface UpdateProtocolSettingsInput {
  allowed_protocol: 'http' | 'https' | 'both';
  apply_scope: 'all' | 'external';
}

@Injectable({ providedIn: 'root' })
export class ProtocolSettingsService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  get(): Observable<ProtocolSettingsOutput> {
    return this.http.get<ProtocolSettingsOutput>(`${this.apiUrl}/protocol-mode`);
  }

  update(data: UpdateProtocolSettingsInput): Observable<ProtocolSettingsOutput> {
    return this.http.put<ProtocolSettingsOutput>(`${this.apiUrl}/protocol-mode`, data);
  }
}
