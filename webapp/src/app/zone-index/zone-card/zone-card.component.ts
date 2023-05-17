import { Component, Input, OnInit } from '@angular/core';
import { humanize_bytes } from 'src/app/core/humanize';
import { ThemeService } from 'src/app/core/services/theme.service';
import { ZoneReportSummary } from 'src/app/olympus-api/zone-report-summary';

@Component({
  selector: 'app-zone-card',
  templateUrl: './zone-card.component.html',
  styleUrls: ['./zone-card.component.scss'],
})
export class ZoneCardComponent implements OnInit {
  public darkTheme: boolean = false;

  @Input() public zone: ZoneReportSummary = new ZoneReportSummary();

  constructor(private theme: ThemeService) {}

  ngOnInit(): void {
    this.theme.isDarkTheme.subscribe((dark) => {
      this.darkTheme = dark;
    });
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
