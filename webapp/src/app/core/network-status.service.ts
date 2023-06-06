import { Injectable } from '@angular/core';
import { Subject, Subscription, fromEvent } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class NetworkStatusService {
  private _online = new Subject<boolean>();
  public online = this._online.asObservable();

  private _subscriptions: Subscription[] = [];
  constructor() {
    this._subscriptions.push(
      fromEvent(window, 'offline').subscribe(() => {
        this._online.next(false);
      })
    );
    this._subscriptions.push(
      fromEvent(window, 'online').subscribe(() => {
        this._online.next(true);
      })
    );
  }

  ngOnDestroy() {
    for (const s of this._subscriptions) {
      s.unsubscribe();
    }
  }
}

@Injectable()
export class ServerNetworkStatusService {
  private _online = new Subject<boolean>();
  public online = this._online.asObservable();
  constructor() {}
}
