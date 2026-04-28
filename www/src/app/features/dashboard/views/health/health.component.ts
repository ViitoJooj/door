import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { forkJoin, Subscription } from 'rxjs';
import { HealthService } from '../../../../core/services/health.service';
import { HealthOverviewData, HealthRouteData, HealthStatus } from '../../dashboard.models';

const ROUTE_LIMIT = 25;

@Component({
  selector: 'app-health',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './health.component.html',
  styleUrl: './health.component.scss'
})
export class HealthComponent implements OnInit, OnDestroy {
  loading = signal(true);
  error = signal<string | null>(null);
  windowMinutes = signal(15);
  overview = signal<HealthOverviewData | null>(null);
  routes = signal<HealthRouteData[]>([]);

  private pollSubscription?: Subscription;

  constructor(private healthService: HealthService) {}

  ngOnInit(): void {
    this.loadHealth();
  }

  ngOnDestroy(): void {
    this.pollSubscription?.unsubscribe();
  }

  setWindow(minutes: number): void {
    if (this.windowMinutes() === minutes) return;
    this.windowMinutes.set(minutes);
    this.loadHealth();
  }

  statusClass(status: HealthStatus): string {
    return `status-badge status-badge--${status}`;
  }

  metricPercent(value: number): string {
    return `${value.toFixed(2)}%`;
  }

  metricLatency(value: number): string {
    return `${value.toFixed(0)} ms`;
  }

  metricRpm(value: number): string {
    return value.toFixed(2);
  }

  private loadHealth(): void {
    this.pollSubscription?.unsubscribe();
    this.loading.set(true);
    this.error.set(null);

    this.pollSubscription = forkJoin({
      overview: this.healthService.getOverview(this.windowMinutes()),
      routes: this.healthService.getRoutes(this.windowMinutes(), ROUTE_LIMIT),
    }).subscribe({
      next: (res) => {
        this.overview.set(res.overview.data);
        this.routes.set(res.routes.data ?? []);
        this.loading.set(false);
      },
      error: (err) => {
        this.error.set(err?.error?.message ?? 'Failed to load health data.');
        this.loading.set(false);
      }
    });
  }
}
