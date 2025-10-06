import type { OnInit } from '@angular/core';
import { Component, inject, signal } from '@angular/core';
import { ApiService } from '../core/services/api.service';

@Component({
  selector: 'app-root-page',
  imports: [],
  templateUrl: './root.component.html',
  styleUrl: './root.component.css',
})
export class RootComponent implements OnInit {
  private readonly apiService = inject(ApiService);

  readonly message = signal('Loading...');

  ngOnInit(): void {
    this.apiService.getRoot().subscribe({
      next: (root) => {
        this.message.set(root.message);
      },
      error: (error) => {
        console.error('Error loading root', error);
        this.message.set('Loading error');
      },
    });
  }
}
