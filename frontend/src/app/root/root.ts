import type { OnInit } from '@angular/core';
import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'root',
  imports: [],
  templateUrl: './root.html',
  styleUrl: './root.css',
})
export class Root implements OnInit {
  private http = inject(HttpClient);
  message = signal<string>('...');

  ngOnInit() {
    this.http
      .get<{ status: string; data: { message: string } }>('/api/')
      .subscribe((res) => this.message.set(res.data.message));
  }
}
