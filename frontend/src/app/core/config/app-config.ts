import type { Provider } from '@angular/core';
import { InjectionToken } from '@angular/core';
import { environment } from '../../../environments/environment';

export interface AppEnvironmentConfig {
  apiUrl: string;
  enableApiLogging: boolean;
  corsAllowedOrigins: string[];
  environmentName: string;
  production: boolean;
}

const normalizeEnvironment = (env: typeof environment): AppEnvironmentConfig => {
  const corsOrigins = env.corsAllowedOrigins
    ? String(env.corsAllowedOrigins)
        .split(',')
        .map((origin) => origin.trim())
        .filter(Boolean)
    : [];

  const apiUrl = (env.apiUrl || '/api').trim().replace(/\/+$/, '');

  return {
    apiUrl: apiUrl || '/api',
    enableApiLogging: Boolean(env.enableApiLogging),
    corsAllowedOrigins: corsOrigins,
    environmentName: env.env ?? (env.production ? 'prod' : 'dev'),
    production: Boolean(env.production),
  };
};

export const APP_ENV_CONFIG = new InjectionToken<AppEnvironmentConfig>('APP_ENV_CONFIG', {
  providedIn: 'root',
  factory: () => normalizeEnvironment(environment),
});

export const createAppEnvironmentConfig = (
  overrides: Partial<AppEnvironmentConfig> = {},
): AppEnvironmentConfig => ({
  ...normalizeEnvironment(environment),
  ...overrides,
});

export const provideAppEnvironmentConfig = (
  config: Partial<AppEnvironmentConfig> | AppEnvironmentConfig,
): Provider => ({
  provide: APP_ENV_CONFIG,
  useValue:
    'production' in config ? (config as AppEnvironmentConfig) : createAppEnvironmentConfig(config),
});
