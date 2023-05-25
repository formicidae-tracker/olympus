import { Component, OnDestroy, OnInit } from '@angular/core';
import { UserSettingsService } from './core/user-settings.service';
import { Subscription, fromEvent } from 'rxjs';
import { MatSnackBar, MatSnackBarRef } from '@angular/material/snack-bar';
import { SnackNetworkOfflineComponent } from './core/snack-network-offline/snack-network-offline.component';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  public darkTheme: boolean = false;

  private _snackOffline?: MatSnackBarRef<any>;
  private _subscriptions: Subscription[] = [];

  constructor(
    private settings: UserSettingsService,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this._subscriptions.push(
      this.settings.isDarkTheme().subscribe((dark) => (this.darkTheme = dark))
    );
    this._subscriptions.push(
      fromEvent(window, 'offline').subscribe(() => {
        this._indicateOffline();
      })
    );
    this._subscriptions.push(
      fromEvent(window, 'online').subscribe(() => {
        this._dismissOffline();
      })
    );
  }

  ngOnDestroy(): void {
    for (const s of this._subscriptions) {
      s.unsubscribe();
    }
  }

  private _indicateOffline() {
    if (this._snackOffline != undefined) {
      return;
    }
    this._snackOffline = this.snackBar.openFromComponent(
      SnackNetworkOfflineComponent
    );
  }

  private _dismissOffline() {
    this._snackOffline?.dismiss();
    this._snackOffline = undefined;
  }
}
