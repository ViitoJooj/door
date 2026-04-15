import { Component, OnInit, OnDestroy, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Subscription, interval, switchMap, startWith } from 'rxjs';
import { RequestLogService } from '../../../../core/services/request-log.service';
import { RequestLog } from '../../dashboard.models';

const POLL_INTERVAL_MS = 3000;

@Component({
  selector: 'app-logs',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './logs.component.html',
  styleUrl: './logs.component.scss'
})
export class LogsComponent implements OnInit, OnDestroy {
  logs = signal<RequestLog[]>([]);
  loading = signal(true);

  private pollSubscription?: Subscription;

  constructor(private requestLogService: RequestLogService) {}

  ngOnInit(): void {
    this.pollSubscription = interval(POLL_INTERVAL_MS).pipe(
      startWith(0),
      switchMap(() => this.requestLogService.getAll())
    ).subscribe({
      next: (res) => {
        const visible = (res.data ?? []).filter(log => log.path !== '/door/api/v1/logs');
        this.logs.set(visible);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      }
    });
  }

  ngOnDestroy(): void {
    this.pollSubscription?.unsubscribe();
  }

  methodClass(method: string): string {
    return `method method--${method.toLowerCase()}`;
  }

  statusClass(code: number): string {
    if (code < 300) return 'status status--2xx';
    if (code < 400) return 'status status--3xx';
    if (code < 500) return 'status status--4xx';
    return 'status status--5xx';
  }
}
