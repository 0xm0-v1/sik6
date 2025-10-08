# HTTP Routing & Data Access Guide

This guide summarises how frontend routes, API calls, and view state are organised in the Angular application.

## Project Layout

```
src/app
  core
    config/app-config.ts        // provides APP_ENV_CONFIG (apiUrl, logging, etc.)
    services/api.service.ts     // centralised HTTP client
    adapters/recipe.adapter.ts  // example view-model mapping
  root
    root.store.ts               // state container used by the root page
  app.routes.ts                 // top-level router configuration
```

### Environment config (APP_ENV_CONFIG)

Instead of importing `environment` directly, inject the `APP_ENV_CONFIG` token. It exposes a normalised object:

```ts
const config = inject(APP_ENV_CONFIG);
console.log(config.apiUrl, config.enableApiLogging);
```

Override it in tests or feature modules with `provideAppEnvironmentConfig({ apiUrl: 'http://localhost:4200/api' })`.

## Creating a New API Call

1. **Add typed models** in `src/app/core/models/api.models.ts`:
   ```ts
   export interface ExampleListData {
     items: ExampleSummary[];
     meta: ApiMeta;
   }
   ```

2. **Expose a service method** using the unified request pipeline:
   ```ts
   // src/app/core/services/api.service.ts
   getExamples(params?: { limit?: number; offset?: number }): Observable<ExampleListData> {
     return this.request<ExampleListData>('GET', '/examples', { params });
   }
   ```
   `request(...)` handles URL building, `HttpParams`, and `ApiEnvelope` extraction.

3. **Adapt the payload** when the UI does not consume raw API shapes. Create a small adapter (see `recipe.adapter.ts`) that maps API fields to view models (Date instances, renamed properties, etc.).

## State & Routing Pattern

The frontend favours a lightweight store per feature that manages API calls and exposes signals. The `RootStore` is an example: it loads the welcome message and recipe list, and the component binds to its signals.

### Steps to add a feature route

1. **Generate the component & store**
   ```bash
   ng generate component examples/list --standalone
   ng generate service examples/examples --flat --skip-tests
   ```
   Replace the generated service with a store similar to `root.store.ts`, injecting `ApiService` and using `takeUntilDestroyed` for subscriptions.

2. **Use the store in the component**
   ```ts
   @Component({
     selector: 'app-examples-page',
     imports: [CommonModule],
     templateUrl: './examples.component.html',
     styleUrl: './examples.component.scss',
   })
   export class ExamplesComponent implements OnInit {
     private readonly store = inject(ExamplesStore);

     readonly items = this.store.items;
     readonly loading = this.store.loading;
     readonly error = this.store.error;

     ngOnInit(): void {
       this.store.initialize();
     }
   }
   ```

3. **Register the route** in `src/app/app.routes.ts`:
   ```ts
   export const routes: Routes = [
     { path: '', component: RootComponent },
     { path: 'examples', loadComponent: () => import('./examples/examples.component').then(m => m.ExamplesComponent) },
   ];
   ```

4. **Template binding**
   ```html
   <section class="examples">
     @if (loading()) {
       <p>Loading examples...</p>
     } @else if (error()) {
       <p class="error">{{ error() }}</p>
     } @else {
       <ul>
         @for (item of items(); track item.id) {
           <li>{{ item.name }}</li>
         } @empty {
           <li><em>No data yet.</em></li>
         }
       </ul>
     }
   </section>
   ```

Signals are functions, so Angular re-renders automatically when their value changes. Use `@for`/`@if` syntax to avoid importing structural directives.

## Local API Development

- The Angular dev server proxies `/api` to the Go backend via `src/proxy.conf.json`.
- Enable extra logging by setting `enableApiLogging` in the environment files or overriding `APP_ENV_CONFIG`.
- When running tests, provide a stubbed `ApiService` or store and use `provideAppEnvironmentConfig` to point at fixture URLs if needed.
