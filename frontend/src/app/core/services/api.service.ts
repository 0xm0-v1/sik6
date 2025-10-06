import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import type { Observable } from 'rxjs';
import { catchError, map, throwError } from 'rxjs';
import type { ApiEnvelope, RootData } from '../models/api.models';
import { ApiError } from '../models/api.models';
import { environment } from '../../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly apiPrefix = environment.apiUrl ?? '/api';

  private get<T>(endpoint: string): Observable<T> {
    return this.http.get<ApiEnvelope<T>>(this.buildUrl(endpoint)).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  private post<T>(endpoint: string, body: unknown): Observable<T> {
    return this.http.post<ApiEnvelope<T>>(this.buildUrl(endpoint), body).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  private put<T>(endpoint: string, body: unknown): Observable<T> {
    return this.http.put<ApiEnvelope<T>>(this.buildUrl(endpoint), body).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  private delete<T>(endpoint: string): Observable<T> {
    return this.http.delete<ApiEnvelope<T>>(this.buildUrl(endpoint)).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  getRoot(): Observable<RootData> {
    return this.get<RootData>('/');
  }

  private buildUrl(endpoint: string): string {
    return `${this.apiPrefix}${endpoint}`;
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

    if (error && typeof error === 'object' && 'status' in error) {
      const httpError = error as { status: number; message?: string };
      const statusCode = httpError.status;
      const errorMessages: Record<number, string> = {
        400: 'Invalid request',
        401: 'Not authorised',
        403: 'Access forbidden',
        404: 'Resource not found',
        500: 'Server error',
        503: 'Service unavailable',
      };

      const message = errorMessages[statusCode] || httpError.message || `HTTP error ${statusCode}`;
      apiError = new ApiError(message, statusCode, error);
    } else if (error instanceof Error) {
      apiError = new ApiError(error.message, undefined, error);
    } else {
      apiError = new ApiError('Unexpected error', undefined, error);
    }

    return throwError(() => apiError);
  }
}
