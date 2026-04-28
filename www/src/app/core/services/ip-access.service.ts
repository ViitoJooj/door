import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { IPAccessListOutput, IPAccessOutput } from '../../features/dashboard/dashboard.models';

export type IPAccessListType = 'whitelist' | 'blacklist';

@Injectable({ providedIn: 'root' })
export class IPAccessService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getAll(type: IPAccessListType): Observable<IPAccessListOutput> {
    return this.http.get<IPAccessListOutput>(`${this.apiUrl}/ip-${type}`);
  }

  create(type: IPAccessListType, ip: string): Observable<IPAccessOutput> {
    return this.http.post<IPAccessOutput>(`${this.apiUrl}/ip-${type}`, { ip });
  }

  update(type: IPAccessListType, id: number, ip: string): Observable<IPAccessOutput> {
    return this.http.put<IPAccessOutput>(`${this.apiUrl}/ip-${type}/${id}`, { ip });
  }

  delete(type: IPAccessListType, id: number): Observable<IPAccessOutput> {
    return this.http.delete<IPAccessOutput>(`${this.apiUrl}/ip-${type}/${id}`);
  }
}
