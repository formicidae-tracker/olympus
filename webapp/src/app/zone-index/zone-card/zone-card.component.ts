import { Component, Input } from '@angular/core';
import { HumanizeService } from 'src/app/core/humanize.service';
import { ZoneReportSummary } from 'src/app/olympus-api/zone-report-summary';

@Component({
  selector: 'app-zone-card',
  templateUrl: './zone-card.component.html',
  styleUrls: ['./zone-card.component.scss'],
})
export class ZoneCardComponent {
  public darkTheme: boolean = false;

  @Input() public zone: ZoneReportSummary = new ZoneReportSummary();
  public subscribed: boolean = false;

  constructor(private humanizer: HumanizeService) {}

  ngOnInit(): void {}

  public usedFraction(): string {
    return this.humanizer.humanizeByteFraction(
      this.zone.tracking?.used_bytes || 0,
      this.zone.tracking?.total_bytes || 0
    );
  }

  public showWarnings(): boolean {
    return this.zone.active_warnings > 0;
  }

  public showEmergencies(): boolean {
    return this.zone.active_emergencies > 0;
  }
}
