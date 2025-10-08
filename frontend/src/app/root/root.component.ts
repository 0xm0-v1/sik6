import type { OnInit } from '@angular/core';
import { Component, inject, signal, PLATFORM_ID } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import type { Recipe } from '../core/models/api.models';
import { ApiService } from '../core/services/api.service';

@Component({
  selector: 'app-root-page',
  // plus besoin de NgIf / NgFor ici
  imports: [],
  templateUrl: './root.component.html',
  styleUrl: './root.component.css',
})
export class RootComponent implements OnInit {
  private readonly apiService = inject(ApiService);
  private readonly platformId = inject(PLATFORM_ID);

  readonly message = signal('Loading...');
  readonly recipes = signal<Recipe[]>([]);
  readonly recipesLoading = signal(true);
  readonly recipesError = signal<string | null>(null);

  ngOnInit(): void {
    if (!isPlatformBrowser(this.platformId)) {
      // Skip API calls during server-side rendering; the browser will load data later.
      this.recipesLoading.set(false);
      this.message.set('Welcome');
      return;
    }

    this.loadRootMessage();
    this.loadRecipes();
  }

  private loadRootMessage(): void {
    this.apiService.getRoot().subscribe({
      next: (root) => this.message.set(root.message),
      error: (error) => {
        console.error('Error loading root', error);
        this.message.set('Loading error');
      },
    });
  }

  private loadRecipes(): void {
    this.recipesLoading.set(true);
    this.recipesError.set(null);

    this.apiService.getRecipes().subscribe({
      next: (payload) => {
        this.recipes.set(payload.recipes);
        this.recipesLoading.set(false);
      },
      error: (error) => {
        console.error('Error loading recipes', error);
        this.recipesLoading.set(false);
        this.recipesError.set(error.message ?? 'Unable to load recipes');
      },
    });
  }
}
