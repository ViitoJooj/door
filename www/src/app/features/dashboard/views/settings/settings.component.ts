import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { EnvService } from '../../../../core/services/env.service';
import { UserAdminService } from '../../../../core/services/user-admin.service';
import { AuthService } from '../../../../core/services/auth.service';
import { RateLimitService } from '../../../../core/services/rate-limit.service';
import { ProtocolSettingsService } from '../../../../core/services/protocol-settings.service';
import { IPAccessListType, IPAccessService } from '../../../../core/services/ip-access.service';
import { SpecialRouteInput, SpecialRoutesService } from '../../../../core/services/special-routes.service';
import { EnvVar, AdminUser, IPAccessData, SpecialRouteData, SpecialRouteType } from '../../dashboard.models';

type SettingsTab = 'env' | 'users' | 'rateLimit' | 'protocol' | 'ipAccess' | 'specialRoutes';

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

  private usersLoaded = false;
  private rateLimitLoaded = false;
  private protocolLoaded = false;
  private whitelistLoaded = false;
  private blacklistLoaded = false;
  private specialLoginLoaded = false;
  private specialRegisterLoaded = false;

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

  rateLimitLoading = signal(false);
  rateLimitSaving = signal(false);
  rateLimitError = signal<string | null>(null);
  rateLimitSuccess = signal<string | null>(null);

  protocolLoading = signal(false);
  protocolSaving = signal(false);
  protocolError = signal<string | null>(null);
  protocolSuccess = signal<string | null>(null);

  activeIPList = signal<IPAccessListType>('whitelist');
  whitelistEntries = signal<IPAccessData[]>([]);
  blacklistEntries = signal<IPAccessData[]>([]);
  ipLoading = signal(false);
  ipCreating = signal(false);
  ipError = signal<string | null>(null);
  ipSuccess = signal<string | null>(null);
  ipNewValue = signal('');
  ipEditingId = signal<number | null>(null);
  ipEditingType = signal<IPAccessListType | null>(null);
  ipEditValue = signal('');
  ipSavingId = signal<number | null>(null);
  ipDeletingId = signal<number | null>(null);

  activeSpecialType = signal<SpecialRouteType>('login');
  loginRoutes = signal<SpecialRouteData[]>([]);
  registerRoutes = signal<SpecialRouteData[]>([]);
  specialLoading = signal(false);
  specialCreating = signal(false);
  specialError = signal<string | null>(null);
  specialSuccess = signal<string | null>(null);
  specialEditingId = signal<number | null>(null);
  specialEditingType = signal<SpecialRouteType | null>(null);
  specialSavingId = signal<number | null>(null);
  specialDeletingId = signal<number | null>(null);
  specialEditValue = signal<SpecialRouteInput | null>(null);

  createForm: FormGroup;
  rateLimitForm: FormGroup;
  protocolForm: FormGroup;
  specialCreateForm: FormGroup;

  constructor(
    private envService: EnvService,
    private userAdminService: UserAdminService,
    private authService: AuthService,
    private rateLimitService: RateLimitService,
    private protocolSettingsService: ProtocolSettingsService,
    private ipAccessService: IPAccessService,
    private specialRoutesService: SpecialRoutesService,
    private fb: FormBuilder
  ) {
    this.isAdmin = this.authService.getCurrentUser()?.role === 'admin';
    this.createForm = this.fb.group({
      username: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      role: ['user', Validators.required],
    });
    this.rateLimitForm = this.fb.group({
      requests_per_second: [10, [Validators.required, Validators.min(0.0001)]],
      burst: [20, [Validators.required, Validators.min(1)]],
      progressive_rate_limit: [true, Validators.required],
    });
    this.protocolForm = this.fb.group({
      allowed_protocol: ['both', Validators.required],
      apply_scope: ['all', Validators.required],
    });
    this.specialCreateForm = this.fb.group({
      path: ['', Validators.required],
      max_distinct_requests: [10, [Validators.required, Validators.min(1)]],
      window_seconds: [60, [Validators.required, Validators.min(1)]],
      ban_seconds: [300, [Validators.required, Validators.min(1)]],
      enabled: [true, Validators.required],
    });
  }

  ngOnInit(): void {
    this.loadEnv();
  }

  setTab(tab: SettingsTab): void {
    this.activeTab.set(tab);
    if (tab === 'users' && this.isAdmin && !this.usersLoaded && !this.usersLoading()) {
      this.loadUsers();
    }
    if (tab === 'rateLimit' && this.isAdmin && !this.rateLimitLoaded && !this.rateLimitLoading()) {
      this.loadRateLimit();
    }
    if (tab === 'protocol' && this.isAdmin && !this.protocolLoaded && !this.protocolLoading()) {
      this.loadProtocolSettings();
    }
    if (tab === 'ipAccess' && !this.ipLoading()) {
      if (!this.whitelistLoaded) this.loadIPList('whitelist');
      if (!this.blacklistLoaded) this.loadIPList('blacklist');
    }
    if (tab === 'specialRoutes' && this.isAdmin && !this.specialLoading()) {
      if (!this.specialLoginLoaded) this.loadSpecialRoutes('login');
      if (!this.specialRegisterLoaded) this.loadSpecialRoutes('register');
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
        this.usersLoaded = true;
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

  loadRateLimit(): void {
    this.rateLimitLoading.set(true);
    this.rateLimitError.set(null);
    this.rateLimitSuccess.set(null);
    this.rateLimitService.get().subscribe({
      next: (res) => {
        this.rateLimitForm.patchValue({
          requests_per_second: res.data.requests_per_second,
          burst: res.data.burst,
          progressive_rate_limit: res.data.progressive_rate_limit,
        });
        this.rateLimitLoading.set(false);
        this.rateLimitLoaded = true;
      },
      error: (err) => {
        this.rateLimitError.set(this.extractError(err, 'Failed to load rate limit settings.'));
        this.rateLimitLoading.set(false);
      }
    });
  }

  saveRateLimit(): void {
    if (this.rateLimitForm.invalid) return;
    this.rateLimitSaving.set(true);
    this.rateLimitError.set(null);
    this.rateLimitSuccess.set(null);
    const payload = this.rateLimitForm.value;
    this.rateLimitService.update({
      requests_per_second: Number(payload.requests_per_second),
      burst: Number(payload.burst),
      progressive_rate_limit: !!payload.progressive_rate_limit,
    }).subscribe({
      next: (res) => {
        this.rateLimitSaving.set(false);
        this.rateLimitSuccess.set(res.message || 'Rate limit updated.');
        this.rateLimitForm.patchValue({
          requests_per_second: res.data.requests_per_second,
          burst: res.data.burst,
          progressive_rate_limit: res.data.progressive_rate_limit,
        });
      },
      error: (err) => {
        this.rateLimitSaving.set(false);
        this.rateLimitError.set(this.extractError(err, 'Failed to update rate limit settings.'));
      }
    });
  }

  loadProtocolSettings(): void {
    this.protocolLoading.set(true);
    this.protocolError.set(null);
    this.protocolSuccess.set(null);
    this.protocolSettingsService.get().subscribe({
      next: (res) => {
        this.protocolForm.patchValue({
          allowed_protocol: res.data.allowed_protocol,
          apply_scope: res.data.apply_scope,
        });
        this.protocolLoading.set(false);
        this.protocolLoaded = true;
      },
      error: (err) => {
        this.protocolError.set(this.extractError(err, 'Failed to load protocol settings.'));
        this.protocolLoading.set(false);
      }
    });
  }

  saveProtocolSettings(): void {
    if (this.protocolForm.invalid) return;
    this.protocolSaving.set(true);
    this.protocolError.set(null);
    this.protocolSuccess.set(null);
    const payload = this.protocolForm.value;
    this.protocolSettingsService.update({
      allowed_protocol: payload.allowed_protocol,
      apply_scope: payload.apply_scope,
    }).subscribe({
      next: (res) => {
        this.protocolSaving.set(false);
        this.protocolSuccess.set(res.message || 'Protocol settings updated.');
        this.protocolForm.patchValue({
          allowed_protocol: res.data.allowed_protocol,
          apply_scope: res.data.apply_scope,
        });
      },
      error: (err) => {
        this.protocolSaving.set(false);
        this.protocolError.set(this.extractError(err, 'Failed to update protocol settings.'));
      }
    });
  }

  currentIPEntries(): IPAccessData[] {
    return this.activeIPList() === 'whitelist' ? this.whitelistEntries() : this.blacklistEntries();
  }

  setIPList(type: IPAccessListType): void {
    this.activeIPList.set(type);
    this.ipError.set(null);
    this.ipSuccess.set(null);
    this.cancelEditIP();
    if (type === 'whitelist' && !this.whitelistLoaded) {
      this.loadIPList('whitelist');
    }
    if (type === 'blacklist' && !this.blacklistLoaded) {
      this.loadIPList('blacklist');
    }
  }

  onIPInput(event: Event): void {
    this.ipNewValue.set((event.target as HTMLInputElement).value);
  }

  loadIPList(type: IPAccessListType): void {
    this.ipLoading.set(true);
    this.ipError.set(null);
    this.ipAccessService.getAll(type).subscribe({
      next: (res) => {
        if (type === 'whitelist') {
          this.whitelistEntries.set(res.data ?? []);
          this.whitelistLoaded = true;
        } else {
          this.blacklistEntries.set(res.data ?? []);
          this.blacklistLoaded = true;
        }
        this.ipLoading.set(false);
      },
      error: (err) => {
        this.ipError.set(this.extractError(err, `Failed to load ${type} entries.`));
        this.ipLoading.set(false);
      }
    });
  }

  createIP(): void {
    const ip = this.ipNewValue().trim();
    if (!ip || this.ipCreating()) return;
    const type = this.activeIPList();
    this.ipCreating.set(true);
    this.ipError.set(null);
    this.ipSuccess.set(null);
    this.ipAccessService.create(type, ip).subscribe({
      next: (res) => {
        this.ipCreating.set(false);
        this.ipNewValue.set('');
        this.ipSuccess.set(res.message || 'IP entry created.');
        this.loadIPList(type);
      },
      error: (err) => {
        this.ipCreating.set(false);
        this.ipError.set(this.extractError(err, `Failed to create ${type} entry.`));
      }
    });
  }

  startEditIP(entry: IPAccessData, type: IPAccessListType): void {
    this.ipEditingId.set(entry.id);
    this.ipEditingType.set(type);
    this.ipEditValue.set(entry.ip);
  }

  onEditIPInput(event: Event): void {
    this.ipEditValue.set((event.target as HTMLInputElement).value);
  }

  cancelEditIP(): void {
    this.ipEditingId.set(null);
    this.ipEditingType.set(null);
    this.ipEditValue.set('');
  }

  saveEditIP(entry: IPAccessData, type: IPAccessListType): void {
    if (this.ipSavingId() !== null || !this.ipEditValue().trim()) return;
    this.ipSavingId.set(entry.id);
    this.ipError.set(null);
    this.ipSuccess.set(null);
    this.ipAccessService.update(type, entry.id, this.ipEditValue().trim()).subscribe({
      next: (res) => {
        this.ipSavingId.set(null);
        this.ipSuccess.set(res.message || 'IP entry updated.');
        this.cancelEditIP();
        this.loadIPList(type);
      },
      error: (err) => {
        this.ipSavingId.set(null);
        this.ipError.set(this.extractError(err, `Failed to update ${type} entry.`));
      }
    });
  }

  deleteIP(entry: IPAccessData, type: IPAccessListType): void {
    this.ipDeletingId.set(entry.id);
    this.ipError.set(null);
    this.ipSuccess.set(null);
    this.ipAccessService.delete(type, entry.id).subscribe({
      next: (res) => {
        this.ipDeletingId.set(null);
        this.ipSuccess.set(res.message || 'IP entry removed.');
        this.loadIPList(type);
      },
      error: (err) => {
        this.ipDeletingId.set(null);
        this.ipError.set(this.extractError(err, `Failed to remove ${type} entry.`));
      }
    });
  }

  currentSpecialRoutes(): SpecialRouteData[] {
    return this.activeSpecialType() === 'login' ? this.loginRoutes() : this.registerRoutes();
  }

  setSpecialType(type: SpecialRouteType): void {
    this.activeSpecialType.set(type);
    this.specialError.set(null);
    this.specialSuccess.set(null);
    this.cancelEditSpecialRoute();
    if (type === 'login' && !this.specialLoginLoaded) {
      this.loadSpecialRoutes('login');
    }
    if (type === 'register' && !this.specialRegisterLoaded) {
      this.loadSpecialRoutes('register');
    }
  }

  loadSpecialRoutes(type: SpecialRouteType): void {
    this.specialLoading.set(true);
    this.specialError.set(null);
    this.specialRoutesService.getAll(type).subscribe({
      next: (res) => {
        if (type === 'login') {
          this.loginRoutes.set(res.data ?? []);
          this.specialLoginLoaded = true;
        } else {
          this.registerRoutes.set(res.data ?? []);
          this.specialRegisterLoaded = true;
        }
        this.specialLoading.set(false);
      },
      error: (err) => {
        this.specialError.set(this.extractError(err, `Failed to load ${type} routes.`));
        this.specialLoading.set(false);
      }
    });
  }

  createSpecialRoute(): void {
    if (this.specialCreateForm.invalid || this.specialCreating()) return;
    this.specialCreating.set(true);
    this.specialError.set(null);
    this.specialSuccess.set(null);
    const type = this.activeSpecialType();
    const payload = this.specialCreateForm.value;
    this.specialRoutesService.create(type, {
      path: payload.path,
      max_distinct_requests: Number(payload.max_distinct_requests),
      window_seconds: Number(payload.window_seconds),
      ban_seconds: Number(payload.ban_seconds),
      enabled: !!payload.enabled,
    }).subscribe({
      next: (res) => {
        this.specialCreating.set(false);
        this.specialSuccess.set(res.message || 'Special route created.');
        this.specialCreateForm.reset({
          path: '',
          max_distinct_requests: 10,
          window_seconds: 60,
          ban_seconds: 300,
          enabled: true,
        });
        this.loadSpecialRoutes(type);
      },
      error: (err) => {
        this.specialCreating.set(false);
        this.specialError.set(this.extractError(err, `Failed to create ${type} route.`));
      }
    });
  }

  startEditSpecialRoute(route: SpecialRouteData): void {
    this.specialEditingId.set(route.id);
    this.specialEditingType.set(route.route_type);
    this.specialEditValue.set({
      path: route.path,
      max_distinct_requests: route.max_distinct_requests,
      window_seconds: route.window_seconds,
      ban_seconds: route.ban_seconds,
      enabled: route.enabled,
    });
  }

  updateSpecialEditField(field: keyof SpecialRouteInput, value: string | number | boolean): void {
    const current = this.specialEditValue();
    if (!current) return;
    this.specialEditValue.set({
      ...current,
      [field]: field === 'enabled'
        ? Boolean(value)
        : (field === 'path' ? String(value) : Number(value)),
    });
  }

  onSpecialEditTextInput(field: 'path', event: Event): void {
    this.updateSpecialEditField(field, (event.target as HTMLInputElement).value);
  }

  onSpecialEditNumberInput(
    field: 'max_distinct_requests' | 'window_seconds' | 'ban_seconds',
    event: Event
  ): void {
    this.updateSpecialEditField(field, Number((event.target as HTMLInputElement).value));
  }

  onSpecialEditEnabledChange(event: Event): void {
    this.updateSpecialEditField('enabled', (event.target as HTMLSelectElement).value === 'true');
  }

  cancelEditSpecialRoute(): void {
    this.specialEditingId.set(null);
    this.specialEditingType.set(null);
    this.specialEditValue.set(null);
  }

  saveEditSpecialRoute(route: SpecialRouteData): void {
    const edit = this.specialEditValue();
    if (!edit || this.specialSavingId() !== null) return;
    this.specialSavingId.set(route.id);
    this.specialError.set(null);
    this.specialSuccess.set(null);
    this.specialRoutesService.update(route.route_type, route.id, edit).subscribe({
      next: (res) => {
        this.specialSavingId.set(null);
        this.specialSuccess.set(res.message || 'Special route updated.');
        this.cancelEditSpecialRoute();
        this.loadSpecialRoutes(route.route_type);
      },
      error: (err) => {
        this.specialSavingId.set(null);
        this.specialError.set(this.extractError(err, 'Failed to update special route.'));
      }
    });
  }

  deleteSpecialRoute(route: SpecialRouteData): void {
    this.specialDeletingId.set(route.id);
    this.specialError.set(null);
    this.specialSuccess.set(null);
    this.specialRoutesService.delete(route.route_type, route.id).subscribe({
      next: (res) => {
        this.specialDeletingId.set(null);
        this.specialSuccess.set(res.message || 'Special route removed.');
        this.loadSpecialRoutes(route.route_type);
      },
      error: (err) => {
        this.specialDeletingId.set(null);
        this.specialError.set(this.extractError(err, 'Failed to remove special route.'));
      }
    });
  }

  private extractError(err: any, fallback: string): string {
    return err?.error?.message ?? fallback;
  }
}
