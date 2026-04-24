import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { EnvService } from '../../../../core/services/env.service';
import { UserAdminService } from '../../../../core/services/user-admin.service';
import { AuthService } from '../../../../core/services/auth.service';
import { EnvVar, AdminUser } from '../../dashboard.models';

type SettingsTab = 'env' | 'users';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './settings.component.html',
  styleUrl: './settings.component.scss'
})
export class SettingsComponent implements OnInit {
  activeTab = signal<SettingsTab>('env');
  isAdmin = false;

  envVars = signal<EnvVar[]>([]);
  envLoading = signal(true);
  envError = signal<string | null>(null);
  editingId = signal<number | null>(null);
  editValue = signal('');
  savingId = signal<number | null>(null);

  users = signal<AdminUser[]>([]);
  usersLoading = signal(false);
  usersError = signal<string | null>(null);
  showCreateUser = signal(false);
  creating = signal(false);
  createError = signal<string | null>(null);
  tempPasswordData = signal<{ username: string; password: string } | null>(null);
  deletingId = signal<number | null>(null);

  createForm: FormGroup;

  constructor(
    private envService: EnvService,
    private userAdminService: UserAdminService,
    private authService: AuthService,
    private fb: FormBuilder
  ) {
    this.isAdmin = this.authService.getCurrentUser()?.role === 'admin';
    this.createForm = this.fb.group({
      username: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      role: ['user', Validators.required],
    });
  }

  ngOnInit(): void {
    this.loadEnv();
  }

  setTab(tab: SettingsTab): void {
    this.activeTab.set(tab);
    if (tab === 'users' && this.users().length === 0 && !this.usersLoading()) {
      this.loadUsers();
    }
  }

  loadEnv(): void {
    this.envLoading.set(true);
    this.envError.set(null);
    this.envService.getAll().subscribe({
      next: (data) => {
        this.envVars.set(data ?? []);
        this.envLoading.set(false);
      },
      error: () => {
        this.envError.set('Failed to load environment variables.');
        this.envLoading.set(false);
      }
    });
  }

  startEdit(env: EnvVar): void {
    this.editingId.set(env.Id);
    this.editValue.set(env.Value);
  }

  cancelEdit(): void {
    this.editingId.set(null);
    this.editValue.set('');
  }

  onEditInput(event: Event): void {
    this.editValue.set((event.target as HTMLInputElement).value);
  }

  saveEdit(env: EnvVar): void {
    if (this.savingId() !== null) return;
    this.savingId.set(env.Id);
    this.envService.update({ ...env, Value: this.editValue() }).subscribe({
      next: (res) => {
        this.envVars.update(vars => vars.map(v => v.Id === env.Id ? res.data : v));
        this.editingId.set(null);
        this.savingId.set(null);
      },
      error: () => {
        this.savingId.set(null);
      }
    });
  }

  loadUsers(): void {
    this.usersLoading.set(true);
    this.usersError.set(null);
    this.userAdminService.getAll().subscribe({
      next: (res) => {
        this.users.set(res.data ?? []);
        this.usersLoading.set(false);
      },
      error: () => {
        this.usersError.set('Failed to load users.');
        this.usersLoading.set(false);
      }
    });
  }

  toggleCreateUser(): void {
    this.showCreateUser.update(v => !v);
    this.createError.set(null);
    this.tempPasswordData.set(null);
    this.createForm.reset({ username: '', email: '', role: 'user' });
  }

  onCreate(): void {
    if (this.createForm.invalid) return;
    this.creating.set(true);
    this.createError.set(null);
    this.tempPasswordData.set(null);

    this.userAdminService.create(this.createForm.value).subscribe({
      next: (res) => {
        this.creating.set(false);
        this.showCreateUser.set(false);
        this.tempPasswordData.set({
          username: res.data.user.username,
          password: res.data.temporary_password,
        });
        this.createForm.reset({ username: '', email: '', role: 'user' });
        this.loadUsers();
      },
      error: (err) => {
        this.creating.set(false);
        this.createError.set(err?.error?.message ?? 'Failed to create user.');
      }
    });
  }

  onDelete(id: number): void {
    this.deletingId.set(id);
    this.userAdminService.delete(id).subscribe({
      next: () => {
        this.deletingId.set(null);
        this.loadUsers();
      },
      error: () => {
        this.deletingId.set(null);
        this.usersError.set('Failed to delete user.');
      }
    });
  }

  dismissTempPassword(): void {
    this.tempPasswordData.set(null);
  }
}
