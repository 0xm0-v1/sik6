import type { OnInit } from '@angular/core';
import { Component, inject, signal } from '@angular/core';
import { ApiService } from '../core/services/api.service';

@Component({
  selector: 'root',
  imports: [],
  templateUrl: './root.html',
  styleUrl: './root.css',
})
export class Root implements OnInit {
  // API service injection (modern Angular syntax)
  private apiService = inject(ApiService);

  // Signal for reactivity (replaces observables in the template)
  message = signal<string>('Loading...');

  ngOnInit(): void {
    // Call to the API service with error handling
    this.apiService.getRoot().subscribe({
      next: (response) => {
        // Direct access to response.data using typing
        if (response.data) {
          this.message.set(response.data.message);
        }
      },
      error: (error) => {
        // Local error handling for UX
        console.error('Error during the loading of Root:', error);
        this.message.set('Loading error');

        // You can add here:
        // - Display a toast/notification
        // - Automatic retry
        // - Redirect to an error page
      },
    });
  }
}
