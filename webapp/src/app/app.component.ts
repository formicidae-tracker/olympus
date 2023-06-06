import { Component, OnDestroy, OnInit } from '@angular/core';
import { UserSettingsService } from './core/user-settings.service';
import { Subscription, fromEvent } from 'rxjs';
import { MatSnackBar, MatSnackBarRef } from '@angular/material/snack-bar';
import { SnackNetworkOfflineComponent } from './core/snack-network-offline/snack-network-offline.component';
import { NetworkStatusService } from './core/network-status.service';

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
    private snackBar: MatSnackBar,
    private networkStatus: NetworkStatusService
  ) {}

  ngOnInit(): void {
    this._subscriptions.push(
      this.settings.isDarkTheme().subscribe((dark) => (this.darkTheme = dark))
    );
    this._subscriptions.push(
      this.networkStatus.online.subscribe((online) => {
        if (online == true) {
          this._dismissOffline();
        } else {
          this._indicateOffline();
        }
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
