# HTTP Routing Guide

This guide explains how frontend routes are organised and how to surface a new backend resource in the Angular application.

## Project Layout

```
src/app
  core
    services/api.service.ts    // centralised HTTP client
    models/api.models.ts       // shared response envelopes
  root                         // example feature module
  app.routes.ts                // top-level router configuration
```

All HTTP calls go through `ApiService`. Components inject it to retrieve data and update their internal `signal` state.

## Expose a New API Call

1. **Add typed models**  
   Extend `src/app/core/models/api.models.ts` with interfaces describing the payload:
   ```ts
   export interface ExampleListData {
     items: ExampleSummary[];
     meta: ApiMeta;
   }

   export type ExampleListResponse = ApiEnvelope<ExampleListData>;
   ```

2. **Expose a service method**  
   Implement a thin wrapper in `src/app/core/services/api.service.ts`:
   ```ts
   getExamples(): Observable<ExampleListData> {
     return this.get<ExampleListData>('/examples');
   }
   ```

   The reusable `get/post/put/delete` helpers already normalise `ApiEnvelope` responses and surface `ApiError` instances when the backend returns a failure.

3. **Create the component**

Generate a standalone component (recommended):
```bash
ng generate component examples/list --standalone
```

Use the `ApiService` inside the component:
```ts
import { Component, OnInit, inject, signal } from '@angular/core';
import { ApiService } from '../core/services/api.service';

@Component({
  selector: 'app-examples',
  templateUrl: './examples.component.html',
  standalone: true,
})
export class ExamplesComponent implements OnInit {
  private readonly api = inject(ApiService);

  readonly items = signal<ExampleListData | null>(null);
  readonly error = signal<string | null>(null);

  ngOnInit(): void {
    this.api.getExamples().subscribe({
      next: (payload) => this.items.set(payload),
      error: (err) => this.error.set(err.message),
    });
  }
}
```

4. **Register the route**

Update `src/app/app.routes.ts`:
```ts
import { Routes } from '@angular/router';
import { RootComponent } from './root/root.component';
import { ExamplesComponent } from './examples/examples.component';

export const routes: Routes = [
  { path: '', component: RootComponent },
  { path: 'examples', component: ExamplesComponent },
];
```

5. **Bind data to the view**

Augment the template to consume the signals exposed by the component and handle loading/error states:
```html
<section class="examples">
  <p *ngIf="loading()">Loading…</p>

  <p *ngIf="!loading() && error()" class="error">
    {{ error() }}
  </p>

  <ul *ngIf="!loading() && !error()">
    <li *ngFor="let item of items(); trackBy: trackItem">
      {{ item.name }}
    </li>
  </ul>
</section>
```

Signals (`loading`, `error`, `items`) are functions, so Angular re-renders automatically when their value changes. Provide a `trackItem` helper in the component for stable list rendering.

During development, the Angular dev server proxies `/api` requests to the Go backend (`src/proxy.conf.json`). As long as the backend exposes the route through the same envelope, the new page will work without extra plumbing.
