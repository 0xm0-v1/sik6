import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import type { Observable } from 'rxjs';
import { catchError, map, throwError } from 'rxjs';
import type { ApiEnvelope, RootResponse } from '../models/api.models';
import { ApiError } from '../models/api.models';

/**
 * Centralized service for all API calls
 *
 * Hybrid architecture:
 * - Private generic methods (get, post, etc.) for reusability
 * - Typed public methods for each endpoint (optimal DX)
 */
@Injectable({
  providedIn: 'root',
})
export class ApiService {
  private http = inject(HttpClient);
  private readonly apiPrefix = '/api';

  /**
   * ============================================
   * GENERIC METHODS (private, reusable)
   * ============================================
   */

  /**
   * Generic GET with error handling
   */
  private get<T>(endpoint: string): Observable<T> {
    const url = `${this.apiPrefix}${endpoint}`;

    return this.http.get<ApiEnvelope<T>>(url).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  /**
   * Generic POST (for future use)
   */
  private post<T>(endpoint: string, body: unknown): Observable<T> {
    const url = `${this.apiPrefix}${endpoint}`;

    return this.http.post<ApiEnvelope<T>>(url, body).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  /**
   * Generic PUT (for future use)
   */
  private put<T>(endpoint: string, body: unknown): Observable<T> {
    const url = `${this.apiPrefix}${endpoint}`;

    return this.http.put<ApiEnvelope<T>>(url, body).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  /**
   * Generic DELETE (for future use)
   */
  private delete<T>(endpoint: string): Observable<T> {
    const url = `${this.apiPrefix}${endpoint}`;

    return this.http.delete<ApiEnvelope<T>>(url).pipe(
      map((envelope) => this.extractData(envelope)),
      catchError((error) => this.handleError(error)),
    );
  }

  /**
   * ============================================
   * PUBLIC METHODS (one per endpoint)
   * ============================================
   */

  /**
   * Retrieves the root route data (/)
   * @returns an Observable containing the root data
   */
  getRoot(): Observable<RootResponse> {
    return this.http
      .get<RootResponse>(`${this.apiPrefix}/`)
      .pipe(catchError((error) => this.handleError(error)));
  }

  /**
   * ============================================
   * UTILITY METHODS (private)
   * ============================================
   */

  /**
   * Extracts data from the API envelope
   * Verifies that the envelope contains data
   */
  private extractData<T>(envelope: ApiEnvelope<T>): T {
    if (envelope.status !== 'ok') {
      throw new ApiError(envelope.error || 'Erreur API inconnue', undefined, envelope);
    }

    if (!envelope.data) {
      throw new ApiError('Réponse API sans données', undefined, envelope);
    }

    return envelope.data;
  }

  /**
   * Centralized HTTP error handling
   * Transforms raw errors into typed ApiErrors
   */
  private handleError(error: unknown): Observable<never> {
    let apiError: ApiError;

    if (error && typeof error === 'object' && 'status' in error) {
      const httpError = error as { status: number; message?: string };
      const statusCode = httpError.status;

      // Custom error messages by HTTP code
      const errorMessages: Record<number, string> = {
        400: 'Requête invalide',
        401: 'Non autorisé',
        403: 'Accès interdit',
        404: 'Ressource introuvable',
        500: 'Erreur serveur',
        503: 'Service indisponible',
      };

      const message = errorMessages[statusCode] || httpError.message || `Erreur HTTP ${statusCode}`;

      apiError = new ApiError(message, statusCode, error);
    } else if (error instanceof Error) {
      apiError = new ApiError(error.message, undefined, error);
    } else {
      apiError = new ApiError('Erreur inconnue', undefined, error);
    }

    return throwError(() => apiError);
  }
}
