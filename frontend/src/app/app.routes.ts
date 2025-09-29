import type { Routes } from '@angular/router';
import { Hello } from './hello/hello';
import { Root } from './root/root';

export const routes: Routes = [
  { path: '', component: Root },
  { path: 'hello', component: Hello },
];
