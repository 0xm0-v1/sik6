/**
 * Generic envelope for all API responses
 * Matches the `response.Envelope` type on the Go backend side
 */
export interface ApiEnvelope<T = unknown> {
  status: string;
  data?: T;
  error?: string;
}

/**
 * Shared metadata emitted by the backend for observability.
 */
export interface ApiMeta {
  component: string;
  type: string;
  time: string;
}

/**
 * Custom API error for better handling
 */
export class ApiError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public originalError?: unknown,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Generic type for data returned by endpoints
 */
export interface ApiData {
  message: string;
  meta: ApiMeta;
}

/**
 * Specific type for root route response (/)
 */
export interface RootData extends ApiData {
  meta: ApiMeta & { type: 'root' };
}

/**
 * Typed envelope types for each endpoint
 */
export type RootResponse = ApiEnvelope<RootData>;

/**
 * Recipe entity mirrors the backend payload.
 */
export interface Recipe {
  id: string;
  recipe_name: string;
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

/**
 * Response payload returned by GET /recipes.
 */
export interface RecipeListData {
  recipes: Recipe[];
  meta: ApiMeta;
}

export type RecipeListResponse = ApiEnvelope<RecipeListData>;
