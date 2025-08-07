import { Injectable } from '@angular/core';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { Observable, timer } from 'rxjs';
import { retryWhen, delayWhen, tap, shareReplay } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class MetricsService {
  private socket$: WebSocketSubject<any> | null = null;
  private readonly WS_URL = 'ws://localhost:8080/ws'; // Replace with your actual WebSocket URL

  private connect(): WebSocketSubject<any> {
    return webSocket({
      url: this.WS_URL,
      openObserver: {
        next: () => console.log('[WebSocket] connected'),
      },
      closeObserver: {
        next: () => {
          console.log('[WebSocket] disconnected');
          this.socket$ = null;
        },
      },
    });
  }

  getMetrics(): Observable<any> {
    if (!this.socket$) {
      this.socket$ = this.connect();
    }

    return this.socket$.pipe(
      retryWhen(errors =>
        errors.pipe(
          tap(err => console.warn('[WebSocket] error, retrying...', err)),
          delayWhen(() => timer(3000))
        )
      ),
      shareReplay(1)
    );
  }
}
