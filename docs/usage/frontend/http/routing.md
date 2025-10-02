# HTTP Routing Guide

This document explains how HTTP routes are organized in the FE ends and how to add new ones.

## Component

1. Generate the component
```bash
    ng generate component [name] [options]
```

2. Component (Hello)
```ts
    import { Component, OnInit, inject, signal } from '@angular/core';
    import { ApiService } from '../core/api.service';

    @Component({
    selector: 'hello',
    imports: [],
    templateUrl: './hello.html',
    styleUrl: './hello.css',
    })
    export class Hello implements OnInit {
    private api = inject(ApiService);
    message = signal<string>('Loading...');

    ngOnInit(): void {
        this.api.getHello().subscribe({
        next: (res) => {
            if (res.data?.message) this.message.set(res.data.message);
        },
        error: (err) => {
            console.error('Hello load error:', err);
            this.message.set('Loading error');
        },
        });
    }
    }
```

3. Register the route at `app.routes.ts`
```ts
    export const routes: Routes = [
    { path: '', component: Root },
    { path: 'hello', component: Hello }, // <- insert route here
    ];
```

4. API typing at `api.models.ts`
```ts
    /**
     * Specific type for hello route response (/)
     */
    export interface HelloData extends ApiData {
    type: 'hello';
    }
```

```ts
    /**
     * Typed envelope types for each endpoint
     */
    export type RootResponse = ApiEnvelope<RootData>;
    export type HelloResponse = ApiEnvelope<HelloData>; 
```

5. Display on a View at `src/app/hello/hello.html`
```html
<p>{{ message() }}</p>
```