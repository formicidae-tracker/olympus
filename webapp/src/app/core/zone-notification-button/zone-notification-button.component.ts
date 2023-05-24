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

  @Input() buttonType: 'icon' | 'flat' | 'fab' = 'icon';

  public disabled: boolean = false;

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
  private _subscriptions: Subscription[] = [];

  constructor(private settings: UserSettingsService) {}

  ngOnInit(): void {
    this.isSolid = this.isSolid !== undefined;
    this._subscriptions.push(
      this.settings.isDarkTheme().subscribe((dark) => {
        this.dark = dark;
      })
    );
    this._subscriptions.push(
      this.settings.getSettings().subscribe((s) => {
        this.disabled = s.subscribeToAll;
      })
    );
  }

  ngOnDestroy(): void {
    this._alarmSubscription?.unsubscribe();
    for (const s of this._subscriptions) {
      s.unsubscribe();
    }
  }

  private _subscribe(): void {
    this._alarmSubscription?.unsubscribe();
    if (this._target.length == 0) {
      return;
    }
    this._alarmSubscription = this.settings
      .hasSubscription(this._target)
      .subscribe((s) => {
        this.subscribed = s;
      });
  }

  public toggleSubscription() {
    if (this._target.length == 0) {
      return;
    }
    if (this.subscribed == true) {
      this.settings.unsubscribeTo(this._target);
    } else {
      this.settings.subscribeTo(this._target);
    }
  }
}
