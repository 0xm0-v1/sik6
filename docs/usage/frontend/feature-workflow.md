# Frontend Feature Workflow

Use this guide when introducing a new feature page (e.g. `Users`) in the Angular app.  
It builds on the existing environment/configuration docs and shows how to wire the store, adapter, route, and tests.

## 1. Plan the data contract

1. Confirm the backend envelope (see `docs/usage/backend/api-pattern.md`) and record the fields you expect.
2. If the payload needs formatting (dates, renamed fields), note the transformation rules.

## 2. Add models & adapters

1. Extend `src/app/core/models/api.models.ts` with API interfaces:
   ```ts
   export interface UserListData {
     users: User[];
     meta: ApiMeta;
   }
   ```
2. Create `src/app/core/adapters/user.adapter.ts`:
   ```ts
   export interface UserViewModel { id: string; name: string; createdAt: Date; }
   export const toUserViewModels = (users: User[]): UserViewModel[] => users.map(...);
   ```
   Adapters keep templates free of backend naming conventions.

## 3. Extend the API client

Add a method to `src/app/core/services/api.service.ts` using the unified request helper:
```ts
getUsers(params?: { limit?: number; offset?: number }): Observable<UserListData> {
  return this.request<UserListData>('GET', '/users', { params });
}
```
The service already injects `APP_ENV_CONFIG` for base URLs.

## 4. Generate the feature shell

1. Create a directory `src/app/users` (or similar).
2. Generate a standalone component:
   ```bash
   ng generate component users/users-page --standalone --inline-style=false --inline-template=false
   ```
3. Create a store:
   ```bash
   ng generate service users/users-store --flat --skip-tests
   ```
   Replace the generated service with a store patterned after `root.store.ts`:
   - Inject `ApiService` and `APP_ENV_CONFIG` if necessary.
   - Maintain signals: `items`, `loading`, `error`.
   - Use `takeUntilDestroyed` and the adapter to map API results.

## 5. Component template & styles

Bind signals via the new control flow syntax:
```html
<section class="users">
  @if (store.loading()) { <p>Loading users...</p> }
  @else if (store.error()) { <p class="error">{{ store.error() }}</p> }
  @else {
    <ul>
      @for (user of store.items(); track user.id) {
        <li>{{ user.name }}</li>
      } @empty {
        <li><em>No users found.</em></li>
      }
    </ul>
  }
</section>
```

## 6. Routing

Update `src/app/app.routes.ts` to expose the page. Prefer lazy loading with `loadComponent`:
```ts
export const routes: Routes = [
  { path: '', component: RootComponent },
  { path: 'users', loadComponent: () => import('./users/users-page.component').then(m => m.UsersPageComponent) },
];
```

## 7. Unit tests

1. Stub the `ApiService` (or provide a fake store) inside the component spec.
2. Use `provideAppEnvironmentConfig` if tests depend on the API base URL.
3. Assert that templates render the expected data & errors.

## 8. End-to-end validation

- Run `npm run lint:check` and `npm run test:ci`.
- Start the stack (`make up`) and verify the new route through the browser.
- Update smoke or integration tests if they must cover the new page.

---

Keep stores small, reuse adapters for presentation logic, and rely on `APP_ENV_CONFIG` to stay environment-agnostic. This mirrors the approach taken by `RootStore`/`RootComponent`.
