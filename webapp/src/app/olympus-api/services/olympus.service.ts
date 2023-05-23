import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { ZoneReportSummary } from '../zone-report-summary';
import { Observable, map } from 'rxjs';
import { ZoneReport } from '../zone-report';
import { ClimateTimeSeries } from '../climate-time-series';

@Injectable({
  providedIn: 'root',
})
export class OlympusService {
  constructor(private httpClient: HttpClient) {}

  getZoneReportSummaries(): Observable<ZoneReportSummary[]> {
    return this.httpClient
      .get<any[]>('/api/zones')
      .pipe(
        map((list) => list.map((item) => ZoneReportSummary.fromPlain(item)))
      );
  }

  getZoneReport(host: string, zone: string): Observable<ZoneReport> {
    return this.httpClient
      .get<any>('/api/host/' + host + '/zone/' + zone)
      .pipe(map((plain) => ZoneReport.fromPlain(plain)));
  }

  getClimateTimeSeries(
    host: string,
    zone: string,
    window: string
  ): Observable<ClimateTimeSeries> {
    return this.httpClient
      .get<any>(
        '/api/host/' + host + '/zone/' + zone + '/climate?window=' + window
      )
      .pipe(map((plain) => ClimateTimeSeries.fromPlain(plain)));
  }
}
