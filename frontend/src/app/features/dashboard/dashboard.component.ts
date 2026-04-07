import { Component } from '@angular/core';
import { AuthService } from '../../core/services/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  template: `
    <div style="padding: 48px; color: #f0f0f0; background: #0e0e10; height: 100vh;">
      <h1 style="font-size: 22px; font-weight: 600; letter-spacing: -0.05em;">door.</h1>
      <p style="margin-top: 16px; color: #505058; font-size: 13px;">Dashboard em construção.</p>
      <button
        (click)="logout()"
        style="margin-top: 24px; background: #f0f0f0; color: #0e0e10; border: none; border-radius: 4px; padding: 8px 16px; font-size: 13px; cursor: pointer;">
        Sair
      </button>
    </div>
  `
})
export class DashboardComponent {
  constructor(private authService: AuthService) {}

  logout(): void {
    this.authService.logout();
  }
}
