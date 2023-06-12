import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { UserSettings } from '../user-settings';
import { LocalStorageService } from './local-storage.service';

export const userSettingsKey = 'userSettings';

@Injectable({
  providedIn: 'root',
})
export class UserSettingsService {
  private settings: UserSettings;
  private dark$: BehaviorSubject<boolean>;
  private alarms$: Map<string, BehaviorSubject<boolean>>;
  // Note we use a Required interface to forbid the access to
  // UserSettings method, they should go through the service. The
  // observable is only useful to get updated data in a component that
  // require everything, i.e. UserSettingsComponent
  private settings$: BehaviorSubject<Required<UserSettings>>;

  constructor(private localStorage: LocalStorageService) {
    this.settings = UserSettings.deserialize(
      this.localStorage.getItem(userSettingsKey) || '{}'
    );
    this.dark$ = new BehaviorSubject<boolean>(this.settings.darkMode);
    this.alarms$ = new Map<string, BehaviorSubject<boolean>>();
    for (const z of this.settings.subscriptions) {
      this.alarms$.set(z, new BehaviorSubject<boolean>(true));
    }
    this.settings$ = new BehaviorSubject<UserSettings>(this.settings);
  }

  public isDarkTheme(): Observable<boolean> {
    return this.dark$.asObservable();
  }

  private _next(): void {
    this.localStorage.setItem(userSettingsKey, this.settings.serialize());
    this.settings$.next(
      new UserSettings(this.settings) as Required<UserSettings>
    );
  }

  public set darkTheme(darkTheme: boolean) {
    if (this.settings.darkMode == darkTheme) {
      return;
    }

    this.settings.darkMode = darkTheme;

    this.dark$.next(darkTheme);
    this._next();
  }

  public set notifyOnWarning(value: boolean) {
    if (this.settings.notifyOnWarning == value) {
      return;
    }
    this.settings.notifyOnWarning = value;
    this._next();
  }

  public set subscribeToAll(value: boolean) {
    if (this.settings.subscribeToAll == value) {
      return;
    }
    this.settings.subscribeToAll = value;
    this._next();
    for (const [zone, subject] of this.alarms$) {
      subject.next(this.settings.hasSubscription(zone));
    }
  }

  public hasSubscription(zone: string): Observable<boolean> {
    let subject = this.alarms$.get(zone);
    if (subject) {
      return subject.asObservable();
    }
    subject = new BehaviorSubject<boolean>(this.settings$.value.subscribeToAll);

    this.alarms$.set(zone, subject);
    return subject.asObservable();
  }

  public getSettings(): Observable<Required<UserSettings>> {
    return this.settings$.asObservable();
  }

  private modifySubscription(zone: string, subscribe: boolean) {
    let skip: boolean;
    if (subscribe) {
      skip = !this.settings.subscribeTo(zone);
    } else {
      skip = !this.settings.unsubscribeFrom(zone);
    }
    if (skip == true) {
      return;
    }

    let subject = this.alarms$.get(zone);
    if (subject == undefined) {
      subject = new BehaviorSubject<boolean>(subscribe);
      this.alarms$.set(zone, subject);
    } else {
      subject.next(subscribe);
    }
    this._next();
  }

  public subscribeTo(zone: string) {
    this.modifySubscription(zone, true);
  }

  public unsubscribeTo(zone: string) {
    this.modifySubscription(zone, false);
  }
}
