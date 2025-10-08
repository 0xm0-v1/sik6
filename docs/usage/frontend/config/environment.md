# Frontend Environment Configuration

Angular code should inject the `APP_ENV_CONFIG` token instead of importing `environment` directly. This keeps runtime values overridable (tests, Storybook, SSR).

```ts
import { inject } from '@angular/core';
import { APP_ENV_CONFIG } from '../core/config/app-config';

const { apiUrl, enableApiLogging } = inject(APP_ENV_CONFIG);
```

## Shape

`AppEnvironmentConfig` contains:

- `apiUrl`: base URL for backend calls (default `/api`). Trailing slashes are stripped.
- `enableApiLogging`: toggles console logging in the API interceptor.
- `corsAllowedOrigins`: array of origins derived from the comma-separated env value.
- `environmentName`: `env` field when provided or `dev`/`prod` fallback.
- `production`: boolean flag mirroring Angular’s environment.

## Providing overrides

Use the helper when a test or feature module needs custom values:

```ts
import { provideAppEnvironmentConfig } from './core/config/app-config';

testBed.configureTestingModule({
  providers: [
    provideAppEnvironmentConfig({
      apiUrl: 'https://fixtures.local/api',
      enableApiLogging: false,
    }),
  ],
});
```

You can also pass a complete `AppEnvironmentConfig` object if you want explicit control.

## SSR considerations

`APP_ENV_CONFIG` is provided at root level, so server-side rendering receives the same defaults. When deploying with environment-specific bundles, update `environment.ts` / `environment.prod.ts`; the token normalises values automatically.

## Relation with ApiService

`ApiService` reads `apiUrl` from the token and relies on its `request(...)` helper for all HTTP verbs. When you override the config, every API call honours the new base URL and logging flag without additional changes.
