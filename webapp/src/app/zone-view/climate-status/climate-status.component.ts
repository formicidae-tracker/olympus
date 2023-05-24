import { Component, Input } from '@angular/core';
import { HumanizeService } from 'src/app/core/humanize.service';
import { ZoneClimateReport } from 'src/app/olympus-api/zone-climate-report';

function formatValue(
  value: number,
  next: number | undefined,
  units: string
): string {
  if (next == undefined) {
    return value.toFixed(1) + units;
  }
  return value.toFixed(1) + '→' + next.toFixed(1) + units;
}

@Component({
  selector: 'app-climate-status',
  templateUrl: './climate-status.component.html',
  styleUrls: ['./climate-status.component.scss'],
})
export class ClimateStatusComponent {
  @Input() climate: ZoneClimateReport = new ZoneClimateReport();
  @Input() now: Date = new Date();
  constructor(private humanizer: HumanizeService) {}

  public columnsToDisplay = ['dim', 'current', 'target'];

  public dataSource() {
    return [
      {
        dim: 'Temperature',
        current: formatValue(this.climate.temperature || 0, undefined, '°C'),
        target: formatValue(
          this.climate.current?.temperature || 0,
          this.climate.current_end?.temperature,
          '°C'
        ),
        next: formatValue(
          this.climate.next?.temperature || 0,
          this.climate.next_end?.temperature,
          '°C'
        ),
      },
      {
        dim: 'Humidity',
        current: formatValue(this.climate.humidity || 0, undefined, '%'),
        target: formatValue(
          this.climate.current?.wind || 0,
          this.climate.current_end?.wind,
          '%'
        ),
        next: formatValue(
          this.climate.next?.wind || 0,
          this.climate.next_end?.wind,
          '%'
        ),
      },

      {
        dim: 'Wind',
        current: undefined,
        target: formatValue(
          this.climate.current?.wind || 0,
          this.climate.current_end?.wind,
          '%'
        ),
        next: formatValue(
          this.climate.next?.wind || 0,
          this.climate.next_end?.wind,
          '%'
        ),
      },
      {
        dim: 'Light',
        current: undefined,
        target: formatValue(
          this.climate.current?.visible_light || 0,
          this.climate.current_end?.visible_light,
          '%'
        ),
        next: formatValue(
          this.climate.next?.visible_light || 0,
          this.climate.next_end?.visible_light,
          '%'
        ),
      },
      {
        dim: 'UV',
        current: undefined,
        target: formatValue(
          this.climate.current?.uv_light || 0,
          this.climate.current_end?.uv_light,
          '%'
        ),
        next: formatValue(
          this.climate.next?.uv_light || 0,
          this.climate.next_end?.uv_light,
          '%'
        ),
      },
    ];
  }

  public trackByIndex(index: number, value: any) {
    return index;
  }

  public formatTarget(): string {
    if (this.climate.next_time == undefined) {
      return 'Target';
    }
    return (
      'Target in ' +
      this.humanizer.humanizeDuration(
        this.climate.next_time.getTime() - this.now.getTime(),
        1
      )
    );
  }
}
