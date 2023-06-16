import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { ZoneReportSummary } from '../zone-report-summary';
import { Observable, map, throwError } from 'rxjs';
import { ZoneReport } from '../zone-report';
import { ClimateTimeSeries } from '../climate-time-series';
import { ServiceLog } from '../service-event';
import { NotificationSettingsUpdate } from '../notification-settings-update';

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

  getLogs(): Observable<ServiceLog[]> {
    return this.httpClient
      .get<any[]>('/api/logs')
      .pipe(
        map((plainList) =>
          plainList.map((plain: any) => ServiceLog.fromPlain(plain))
        )
      );
  }

  getVersion(): Observable<string> {
    return this.httpClient
      .get<any>('/api/version')
      .pipe(map((plain: any) => plain.version || 'undefined'));
  }

  getPushServerPublicKey(): Observable<string> {
    return this.httpClient
      .get<any>('/api/notifications/key')
      .pipe(map((plain: any) => plain.serverPublicKey || ''));
  }

  registerPushSubscription(s: PushSubscription): Observable<void> {
    return this.httpClient
      .post<PushSubscription>('/api/notifications', s)
      .pipe(map(() => void 0));
  }

  updateNotificationSettings(u: NotificationSettingsUpdate): Observable<void> {
    return this.httpClient
      .post<string>('/api/notifications/settings', u.serialize(), {
        headers: { 'Content-Type': 'application/json' },
      })
      .pipe(map(() => void 0));
  }
}
