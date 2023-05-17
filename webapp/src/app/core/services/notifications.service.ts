import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class NotificationsSettings {
  public zones: Set<string> = new Set<string>();

  static fromJSON(jsondata: string): NotificationsSettings {
    let plain = JSON.parse(jsondata);
    let res = new NotificationsSettings();
    res.zones = new Set<string>(plain.zones || []);
    return res;
  }

  public subscribed(zoneIdentifier: string): boolean {
    return this.zones.has(zoneIdentifier);
  }

  toJSON(): string {
    return this.toJSON();
  }
}

const key = 'notificationSettings';

export class NotificationsService {
  private _subscriptions: BehaviorSubject<NotificationsSettings>;
  public subscriptions: Observable<NotificationsSettings>;

  constructor() {
    let stored = NotificationsSettings.fromJSON(
      localStorage.getItem(key) || '{}'
    );
    this._subscriptions = new BehaviorSubject<NotificationsSettings>(stored);
    this.subscriptions = this._subscriptions.asObservable();
  }

  public subscribe(zoneIdentifier: string) {
    let settings = this._subscriptions.value;
    if (settings.subscribed(zoneIdentifier)) {
      return;
    }
    settings.zones.add(zoneIdentifier);
    this._subscriptions.next(settings);
    localStorage.setItem(key, settings.toJSON());
  }

  public unsubscribe(zoneIdentifier: string): void {
    let settings = this._subscriptions.value;
    if (settings.subscribed(zoneIdentifier) == false) {
      return;
    }
    settings.zones.delete(zoneIdentifier);
    this._subscriptions.next(settings);
    localStorage.setItem(key, settings.toJSON());
  }
}
