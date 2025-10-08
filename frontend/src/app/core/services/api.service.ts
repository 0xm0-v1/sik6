import {
  HttpClient,
  HttpErrorResponse,
  HttpParams,
  type HttpParameterCodec,
} from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import type { Observable } from 'rxjs';
import { catchError, map, throwError } from 'rxjs';
import type { ApiEnvelope, RecipeListData, RootData } from '../models/api.models';
import { ApiError } from '../models/api.models';
import { APP_ENV_CONFIG } from '../config/app-config';

type HttpMethod = 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly envConfig = inject(APP_ENV_CONFIG);

  getRoot(): Observable<RootData> {
    return this.request<RootData>('GET', '/');
  }

  getRecipes(params?: { limit?: number; offset?: number }): Observable<RecipeListData> {
    return this.request<RecipeListData>('GET', '/recipes', { params });
  }

  private request<T>(
    method: HttpMethod,
    endpoint: string,
    options: {
      body?: unknown;
      params?: Record<string, string | number | boolean | undefined | null>;
    } = {},
  ): Observable<T> {
    const url = this.buildUrl(endpoint);
    const requestOptions: { body?: unknown; params?: HttpParams } = {};

    if (options.body !== undefined) {
      requestOptions.body = options.body;
    }

    if (options.params) {
      requestOptions.params = this.buildParams(options.params);
    }

    return this.http.request<ApiEnvelope<T>>(method, url, requestOptions).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  private buildParams(
    params: Record<string, string | number | boolean | undefined | null>,
  ): HttpParams {
    let httpParams = new HttpParams({ encoder: new StrictHttpParamEncoder() });

    for (const [key, value] of Object.entries(params)) {
      if (value === undefined || value === null) {
        continue;
      }
      httpParams = httpParams.set(key, String(value));
    }

    return httpParams;
  }

  private buildUrl(endpoint: string): string {
    const normalizedEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    return `${this.envConfig.apiUrl}${normalizedEndpoint}`;
  }

  private extractData<T>(envelope: ApiEnvelope<T>): T {
    if (envelope.status !== 'ok') {
      throw new ApiError(envelope.error || 'Unknown API error', undefined, envelope);
    }

    if (envelope.data == null) {
      throw new ApiError('API response missing data', undefined, envelope);
    }

    return envelope.data;
  }

  private handleError(error: unknown): Observable<never> {
    let apiError: ApiError;

    if (error instanceof HttpErrorResponse) {
      const statusCode = error.status;
      const errorMessages: Record<number, string> = {
        400: 'Invalid request',
        401: 'Not authorised',
        403: 'Access forbidden',
        404: 'Resource not found',
        500: 'Server error',
        503: 'Service unavailable',
      };

      const message = errorMessages[statusCode] || error.message || `HTTP error ${statusCode}`;
      apiError = new ApiError(message, statusCode, error.error ?? error);
    } else if (error instanceof Error) {
      apiError = new ApiError(error.message, undefined, error);
    } else {
      apiError = new ApiError('Unexpected error', undefined, error);
    }

    return throwError(() => apiError);
  }
}

/**
 * HttpParams encoder that avoids automatic plus/space conversion.
 */
class StrictHttpParamEncoder implements HttpParameterCodec {
  encodeKey(key: string): string {
    return encodeURIComponent(key);
  }

  encodeValue(value: string): string {
    return encodeURIComponent(value);
  }

  decodeKey(key: string): string {
    return decodeURIComponent(key);
  }

  decodeValue(value: string): string {
    return decodeURIComponent(value);
  }
}
