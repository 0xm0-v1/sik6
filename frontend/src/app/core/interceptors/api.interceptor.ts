/* eslint-disable no-console */
import { inject } from '@angular/core';
import type { HttpInterceptorFn } from '@angular/common/http';
import { tap } from 'rxjs';
import { APP_ENV_CONFIG } from '../config/app-config';

/**
 * Global HTTP interceptor for:
 * - Logging all requests/responses (depending on configuration)
 * - Monitoring errors
 * - Adding global headers if needed (future)
 *
 * Note: Use Angular's new functional format (no classes)
 * Configure via environment.enableApiLogging
 */
export const apiInterceptor: HttpInterceptorFn = (req, next) => {
  const { enableApiLogging } = inject(APP_ENV_CONFIG);
  const isBrowser = typeof window !== 'undefined';

  if (isBrowser && enableApiLogging) {
    console.log(`[API] ${req.method} ${req.url}`);
  }

  return next(req).pipe(
    tap({
      next: (event) => {
        if (isBrowser && enableApiLogging) {
          if ('status' in event) {
            console.log(`[API] OK ${req.method} ${req.url} - ${event.status}`);
          }
        }
      },
      error: (error) => {
        console.error(`[API] ERR ${req.method} ${req.url}`, error);
      },
    }),
  );
};
