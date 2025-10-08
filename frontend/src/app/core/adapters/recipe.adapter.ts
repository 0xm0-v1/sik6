import type { Recipe } from '../models/api.models';

export interface RecipeViewModel {
  id: string;
  name: string;
  createdAt: Date;
  updatedAt: Date;
  deletedAt: Date | null;
}

export const toRecipeViewModel = (recipe: Recipe): RecipeViewModel => ({
  id: recipe.id,
  name: recipe.recipe_name,
  createdAt: new Date(recipe.created_at),
  updatedAt: new Date(recipe.updated_at),
  deletedAt: recipe.deleted_at ? new Date(recipe.deleted_at) : null,
});

export const toRecipeViewModels = (recipes: Recipe[]): RecipeViewModel[] =>
  recipes.map((recipe) => toRecipeViewModel(recipe));
