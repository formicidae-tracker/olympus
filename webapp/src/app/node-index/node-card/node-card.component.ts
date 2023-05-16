import { Component, Input } from '@angular/core';
import { humanize_bytes } from 'src/app/core/humanize';
import { ZoneReportSummary } from 'src/app/olympus-api/zone-report-summary';

@Component({
  selector: 'app-node-card',
  templateUrl: './node-card.component.html',
  styleUrls: ['./node-card.component.scss'],
})
export class NodeCardComponent {
  @Input() public zone: ZoneReportSummary = new ZoneReportSummary();

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
