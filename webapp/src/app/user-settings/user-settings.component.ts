import {
  Component,
  ElementRef,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { UserSettingsService } from '../core/services/user-settings.service';
import { OlympusService } from '../olympus-api/services/olympus.service';
import { Observable, Subscription, map, of } from 'rxjs';
import { MatChipInputEvent } from '@angular/material/chips';
import { MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { UserSettings } from '../core/user-settings';

@Component({
  selector: 'app-user-settings',
  templateUrl: './user-settings.component.html',
  styleUrls: ['./user-settings.component.scss'],
})
export class UserSettingsComponent implements OnInit, OnDestroy {
  //Note we use a Required interface to disallow calling UserSetting
  //logic.
  public settings: Required<UserSettings> = new UserSettings();
  private _availableZones: string[] = [];
  private _subscriptions: Subscription[] = [];

  public get subscribeToAll(): boolean {
    return this.settings.subscribeToAll;
  }
  public set subscribeToAll(value: boolean) {
    this.settingsService.subscribeToAll = value;
  }

  public get darkTheme(): boolean {
    return this.settings.darkMode;
  }
  public set darkTheme(value: boolean) {
    this.settingsService.darkTheme = value;
  }

  public get notifyOnWarnings(): boolean {
    return this.settings.notifyOnWarning;
  }
  public set notifyOnWarnings(value: boolean) {
    this.settingsService.notifyOnWarning = value;
  }

  @ViewChild('zoneInput') zoneInput!: ElementRef<HTMLInputElement>;

  constructor(
    private settingsService: UserSettingsService,
    private olympus: OlympusService
  ) {}

  ngOnInit(): void {
    this._subscriptions.push(
      this.settingsService.getSettings().subscribe((settings) => {
        this.settings = settings;
      })
    );
    this.olympus.getZoneReportSummaries().subscribe((zones) => {
      this._availableZones = zones.map((zone) => zone.identifier);
    });
  }

  ngOnDestroy(): void {
    for (const s of this._subscriptions) {
      s.unsubscribe();
    }
  }

  public subscribeTo(event: MatChipInputEvent): void {
    const zone: string = (event.value || '').trim();
    if (zone) {
      this.settings.subscribeTo(zone);
    }
    event.chipInput!.clear();
  }

  public selected(event: MatAutocompleteSelectedEvent): void {
    this.settingsService.subscribeTo(event.option.viewValue);
    this.zoneInput.nativeElement.value = '';
  }

  public unsubscribeFrom(zoneIdentifier: string) {
    this.settingsService.unsubscribeTo(zoneIdentifier);
  }

  public zones(): Observable<string[]> {
    return of(this._availableZones).pipe(
      map((zones) => {
        return Array.from(new Set<string>(zones)).filter(
          (zone) => !this.settings.hasSubscription(zone)
        );
      })
    );
  }
}
