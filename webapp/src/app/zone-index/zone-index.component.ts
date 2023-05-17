import { Component, OnInit } from '@angular/core';
import { ZoneReportSummary } from '../olympus-api/zone-report-summary';
import { OlympusService } from '../olympus-api/services/olympus.service';

@Component({
  selector: 'app-zone-index',
  templateUrl: './zone-index.component.html',
  styleUrls: ['./zone-index.component.scss'],
})
export class ZoneIndexComponent implements OnInit {
  public zones: ZoneReportSummary[] = [];

  constructor(private olympus: OlympusService) {}

  ngOnInit(): void {
    this.olympus
      .getZoneReportSummaries()
      .subscribe((zones) => (this.zones = zones));
  }
}
