import type { OnInit } from '@angular/core';
import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RootStore } from './root.store';

@Component({
  selector: 'app-root-page',
  imports: [CommonModule],
  templateUrl: './root.component.html',
  styleUrl: './root.component.css',
})
export class RootComponent implements OnInit {
  private readonly store = inject(RootStore);

  readonly message = this.store.message;
  readonly recipes = this.store.recipes;
  readonly recipesLoading = this.store.recipesLoading;
  readonly recipesError = this.store.recipesError;

  ngOnInit(): void {
    this.store.initialize();
  }

  refreshRecipes(): void {
    this.store.refreshRecipes();
  }
}
