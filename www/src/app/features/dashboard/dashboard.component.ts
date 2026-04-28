import { Component, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../core/services/auth.service';
import { ApplicationsComponent } from './views/applications/applications.component';
import { HealthComponent } from './views/health/health.component';
import { LogsComponent } from './views/logs/logs.component';
import { SettingsComponent } from './views/settings/settings.component';

type ActiveView = 'applications' | 'health' | 'logs' | 'settings';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, ApplicationsComponent, HealthComponent, LogsComponent, SettingsComponent],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss'
})
export class DashboardComponent {
  activeView = signal<ActiveView>('applications');

  constructor(private authService: AuthService) {}

  setView(view: ActiveView): void {
    this.activeView.set(view);
  }

  logout(): void {
    this.authService.logout();
  }
}
