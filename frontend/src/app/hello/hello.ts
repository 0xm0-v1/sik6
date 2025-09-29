import type { OnInit } from '@angular/core';
import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'hello',
  imports: [],
  templateUrl: './hello.html',
  styleUrl: './hello.css',
})
export class Hello implements OnInit {
  private http = inject(HttpClient);
  message = signal<string>('...');

  ngOnInit() {
    this.http
      .get<{ status: string; data: { message: string } }>('/api/hello')
      .subscribe((res) => this.message.set(res.data.message));
  }
}
