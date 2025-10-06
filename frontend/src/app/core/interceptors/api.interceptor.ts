import type { HttpInterceptorFn } from '@angular/common/http';
import { tap } from 'rxjs';
import { environment } from '../../../environments/environment';

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
  if (typeof window !== 'undefined' && environment.enableApiLogging) {
    console.log(`[API] ${req.method} ${req.url}`);
  }

  return next(req).pipe(
    tap({
      next: (event) => {
        if (typeof window !== 'undefined' && environment.enableApiLogging) {
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
