import { Component, Input } from '@angular/core';
import { HumanizeService } from 'src/app/core/services/humanize.service';
import { TrackingInfo } from 'src/app/olympus-api/tracking-info';

@Component({
  selector: 'app-tracking-status',
  templateUrl: './tracking-status.component.html',
  styleUrls: ['./tracking-status.component.scss'],
})
export class TrackingStatusComponent {
  @Input() tracking: TrackingInfo = new TrackingInfo();
  @Input() now: Date = new Date();
  constructor(private humanizer: HumanizeService) {}

  public duration(): string {
    return this.humanizer.humanizeDuration(
      this.now.getTime() - this.tracking.since.getTime(),
      1
    );
  }

  public formatDiskSpace(): string {
    return this.humanizer.humanizeByteFraction(
      this.tracking.used_bytes,
      this.tracking.total_bytes
    );
  }

  public formatFillRate(): string {
    return this.humanizer.humanizeBytes(this.tracking.bytes_per_second) + '/s';
  }

  public formatFillETA(): string {
    return this.humanizer.humanizeDuration(this.tracking.filledUpEta(), 1);
  }
}
