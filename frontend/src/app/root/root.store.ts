/* eslint-disable no-console */
import { DestroyRef, Injectable, PLATFORM_ID, computed, inject, signal } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { ApiError } from '../core/models/api.models';
import { ApiService } from '../core/services/api.service';
import { toRecipeViewModels, type RecipeViewModel } from '../core/adapters/recipe.adapter';

@Injectable({
  providedIn: 'root',
})
export class RootStore {
  private readonly apiService = inject(ApiService);
  private readonly platformId = inject(PLATFORM_ID);
  private readonly destroyRef = inject(DestroyRef);

  readonly message = signal<string>('Loading...');
  readonly recipes = signal<RecipeViewModel[]>([]);
  readonly recipesLoading = signal<boolean>(false);
  readonly recipesError = signal<string | null>(null);
  readonly hasRecipes = computed(() => this.recipes().length > 0);

  initialize(): void {
    if (!isPlatformBrowser(this.platformId)) {
      this.message.set('Welcome');
      this.recipesLoading.set(false);
      return;
    }

    this.fetchRootMessage();
    this.fetchRecipes();
  }

  refreshRecipes(): void {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }
    this.fetchRecipes();
  }

  private fetchRootMessage(): void {
    this.apiService
      .getRoot()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe({
        next: (root) => this.message.set(root.message),
        error: (error) => {
          console.error('Error loading root', error);
          this.message.set('Loading error');
        },
      });
  }

  private fetchRecipes(): void {
    this.recipesLoading.set(true);
    this.recipesError.set(null);

    this.apiService
      .getRecipes()
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe({
        next: (payload) => {
          this.recipes.set(toRecipeViewModels(payload.recipes));
          this.recipesLoading.set(false);
        },
        error: (error: unknown) => {
          console.error('Error loading recipes', error);
          this.recipesLoading.set(false);
          const message = error instanceof ApiError ? error.message : 'Unable to load recipes';
          this.recipesError.set(message);
        },
      });
  }
}
