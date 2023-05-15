import { Component, Input } from '@angular/core';
import { ZoneReportSummary } from 'src/app/olympus-api/zone-report-summary';

@Component({
  selector: 'app-node-card',
  templateUrl: './node-card.component.html',
  styleUrls: ['./node-card.component.scss'],
})
export class NodeCardComponent {
  @Input() public zone: ZoneReportSummary = new ZoneReportSummary();
}
