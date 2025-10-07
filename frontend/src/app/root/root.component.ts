import type { OnInit } from '@angular/core';
import { Component, inject, signal } from '@angular/core';
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

  readonly message = signal('Loading...');
  readonly recipes = signal<Recipe[]>([]);
  readonly recipesLoading = signal(true);
  readonly recipesError = signal<string | null>(null);

  ngOnInit(): void {
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
