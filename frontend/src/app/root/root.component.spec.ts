import { PLATFORM_ID, provideZonelessChangeDetection } from '@angular/core';
import type { ComponentFixture } from '@angular/core/testing';
import { TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { RootComponent } from './root.component';
import { ApiService } from '../core/services/api.service';

describe('RootComponent', () => {
  let component: RootComponent;
  let fixture: ComponentFixture<RootComponent>;

  class ApiServiceStub {
    getRoot() {
      return of({
        message: 'Hello',
        meta: {
          component: 'api',
          type: 'root' as const,
          time: new Date().toISOString(),
        },
      });
    }

    getRecipes() {
      return of({
        recipes: [],
        meta: {
          component: 'api',
          type: 'recipes:list',
          time: new Date().toISOString(),
        },
      });
    }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RootComponent],
      providers: [
        provideZonelessChangeDetection(),
        { provide: ApiService, useClass: ApiServiceStub },
        { provide: PLATFORM_ID, useValue: 'browser' },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(RootComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should render the root message', () => {
    const compiled = fixture.nativeElement as HTMLElement;
    expect(compiled.textContent).toContain('Hello');
  });
});
