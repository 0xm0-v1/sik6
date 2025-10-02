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
  component: string;
  type: string;
  time: string;
}

/**
 * Specific type for root route response (/)
 */
export interface RootData extends ApiData {
  type: 'root';
}

/**
 * Typed envelope types for each endpoint
 */
export type RootResponse = ApiEnvelope<RootData>;
