import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { UserSettingsService } from '../user-settings.service';

@Component({
  selector: 'app-zone-notification-button',
  templateUrl: './zone-notification-button.component.html',
  styleUrls: ['./zone-notification-button.component.scss'],
})
export class ZoneNotificationButtonComponent implements OnInit, OnDestroy {
  @Input('solid') isSolid: boolean | string | undefined;

  @Input() buttonType: 'icon' | 'flat' = 'icon';

  @Input()
  set target(value: string) {
    this._target = value;
    this._subscribe();
  }

  get target(): string {
    return this._target;
  }

  public subscribed: boolean = false;
  public dark: boolean = false;

  private _target = '';
  private _alarmSubscription?: Subscription;
  private _themeSubscription?: Subscription;

  constructor(private settings: UserSettingsService) {}

  ngOnInit(): void {
    this.isSolid = this.isSolid !== undefined;
    this._themeSubscription = this.settings.isDarkTheme().subscribe((dark) => {
      this.dark = dark;
    });
  }

  ngOnDestroy(): void {
    this._alarmSubscription?.unsubscribe();
    this._themeSubscription?.unsubscribe();
  }

  private _subscribe(): void {
    this._alarmSubscription?.unsubscribe();
    if (this._target.length == 0) {
      return;
    }
    this._alarmSubscription = this.settings
      .isSubscribedToAlarmFrom(this._target)
      .subscribe((s) => {
        this.subscribed = s;
      });
  }

  public toggleSubscription() {
    if (this._target.length == 0) {
      return;
    }
    if (this.subscribed == true) {
      this.settings.unsubscribeFromAlarmFrom(this._target);
    } else {
      this.settings.subscribeToAlarmFrom(this._target);
    }
  }
}
