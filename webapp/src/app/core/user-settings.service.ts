import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { UserSettings } from './user-settings';

export const userSettingsKey = 'userSettings';

@Injectable({
  providedIn: 'root',
})
export class UserSettingsService {
  private userSettings: UserSettings;
  private darkTheme: BehaviorSubject<boolean>;
  private alarmSubscription: Map<string, BehaviorSubject<boolean>>;

  constructor() {
    this.userSettings = UserSettings.deserialize(
      localStorage.getItem(userSettingsKey) || '{}'
    );
    this.darkTheme = new BehaviorSubject<boolean>(this.userSettings.darkMode);
    this.alarmSubscription = new Map<string, BehaviorSubject<boolean>>();
    for (const z of this.userSettings.alarmSubscriptions) {
      this.alarmSubscription.set(z, new BehaviorSubject<boolean>(true));
    }
  }

  public isDarkTheme(): Observable<boolean> {
    return this.darkTheme.asObservable();
  }

  private saveToLocalStorage(): void {
    localStorage.setItem(userSettingsKey, this.userSettings.serialize());
  }

  public setDarkTheme(darkTheme: boolean): void {
    if (this.userSettings.darkMode == darkTheme) {
      return;
    }
    this.userSettings.darkMode = darkTheme;
    this.darkTheme.next(darkTheme);
    this.saveToLocalStorage();
  }

  public isSubscribedToAlarmFrom(zone: string): Observable<boolean> {
    let subject = this.alarmSubscription.get(zone);
    if (subject) {
      return subject.asObservable();
    }
    subject = new BehaviorSubject<boolean>(false);
    this.alarmSubscription.set(zone, subject);
    return subject.asObservable();
  }

  private modifyAlarmSubscription(zone: string, subscribe: boolean) {
    let skip: boolean;
    if (subscribe) {
      skip = !this.userSettings.subscribeToAlarmFrom(zone);
    } else {
      skip = !this.userSettings.unsubscribeFromAlarmFrom(zone);
    }

    if (skip) {
      return;
    }

    let subject = this.alarmSubscription.get(zone);
    if (subject == undefined) {
      subject = new BehaviorSubject<boolean>(subscribe);
      this.alarmSubscription.set(zone, subject);
    } else {
      subject.next(subscribe);
    }

    this.saveToLocalStorage();
  }

  public subscribeToAlarmFrom(zone: string) {
    this.modifyAlarmSubscription(zone, true);
  }

  public unsubscribeFromAlarmFrom(zone: string) {
    this.modifyAlarmSubscription(zone, false);
  }
}
