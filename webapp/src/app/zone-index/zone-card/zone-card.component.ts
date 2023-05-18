import { Component, Input, OnInit } from '@angular/core';
import { humanize_bytes } from 'src/app/core/humanize';
import { UserSettingsService } from 'src/app/core/user-settings.service';
import { ZoneReportSummary } from 'src/app/olympus-api/zone-report-summary';

@Component({
  selector: 'app-zone-card',
  templateUrl: './zone-card.component.html',
  styleUrls: ['./zone-card.component.scss'],
})
export class ZoneCardComponent implements OnInit {
  public darkTheme: boolean = false;

  @Input() public zone: ZoneReportSummary = new ZoneReportSummary();
  public subscribed: boolean = false;

  constructor(private settingsService: UserSettingsService) {}

  ngOnInit(): void {
    this.settingsService.isDarkTheme().subscribe((dark) => {
      this.darkTheme = dark;
    });
    this.settingsService
      .isSubscribedToAlarmFrom(this.zone.identifier())
      .subscribe((subscribed) => {
        this.subscribed = subscribed;
      });
  }

  public setAlarmSubscription(subscribed: boolean): void {
    if (subscribed) {
      this.settingsService.subscribeToAlarmFrom(this.zone.identifier());
    } else {
      this.settingsService.unsubscribeFromAlarmFrom(this.zone.identifier());
    }
  }

  public fill_rate(): string {
    return humanize_bytes(
      this.zone.tracking ? this.zone.tracking.bytes_per_second : 0,
      'B/s'
    );
  }

  public used_space(): string {
    return humanize_bytes(
      this.zone.tracking
        ? this.zone.tracking.total_bytes - this.zone.tracking.free_bytes
        : 0
    );
  }

  public total_space(): string {
    return humanize_bytes(
      this.zone.tracking ? this.zone.tracking.total_bytes : 0
    );
  }
}
