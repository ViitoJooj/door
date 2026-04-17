import { Component, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { ApplicationService } from '../../../../core/services/application.service';
import { Application } from '../../dashboard.models';

@Component({
  selector: 'app-applications',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './applications.component.html',
  styleUrl: './applications.component.scss'
})
export class ApplicationsComponent implements OnInit {
  applications = signal<Application[]>([]);
  loading = signal(true);
  errorMessage = signal<string | null>(null);
  showCreateForm = signal(false);
  creating = signal(false);
  createError = signal<string | null>(null);
  deletingId = signal<number | null>(null);

  createForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private applicationService: ApplicationService
  ) {
    this.createForm = this.fb.group({
      url: ['', Validators.required],
      country: ['', Validators.required]
    });
  }

  ngOnInit(): void {
    this.loadApplications();
  }

  loadApplications(): void {
    this.loading.set(true);
    this.errorMessage.set(null);

    this.applicationService.getAll().subscribe({
      next: (data) => {
        this.applications.set(data ?? []);
        this.loading.set(false);
      },
      error: () => {
        this.errorMessage.set('Erro ao carregar aplicações.');
        this.loading.set(false);
      }
    });
  }

  toggleCreateForm(): void {
    this.showCreateForm.update(v => !v);
    this.createError.set(null);
    this.createForm.reset();
  }

  onCreate(): void {
    if (this.createForm.invalid) return;

    this.creating.set(true);
    this.createError.set(null);

    this.applicationService.create(this.createForm.value).subscribe({
      next: () => {
        this.creating.set(false);
        this.showCreateForm.set(false);
        this.createForm.reset();
        this.loadApplications();
      },
      error: () => {
        this.creating.set(false);
        this.createError.set('Erro ao criar aplicação. Verifique os dados e tente novamente.');
      }
    });
  }

  onDelete(id: number): void {
    this.deletingId.set(id);

    this.applicationService.delete(id).subscribe({
      next: () => {
        this.deletingId.set(null);
        this.loadApplications();
      },
      error: () => {
        this.deletingId.set(null);
        this.errorMessage.set('Erro ao remover aplicação.');
      }
    });
  }
}
