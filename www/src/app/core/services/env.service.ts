import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { EnvVar, EnvOutput } from '../../features/dashboard/dashboard.models';

@Injectable({ providedIn: 'root' })
export class EnvService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(): Observable<EnvVar[]> {
    return this.http.get<EnvVar[]>(`${this.apiUrl}/env/`);
  }

  update(env: EnvVar): Observable<EnvOutput> {
    return this.http.put<EnvOutput>(`${this.apiUrl}/env/${env.Id}`, {
      id: env.Id,
      name: env.Name,
      value: env.Value,
    });
  }
}
