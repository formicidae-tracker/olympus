import { Injectable } from '@angular/core';
import { ZoneClimateReport,ZoneClimateReportAdapter } from '@models/zone-climate-report';
import { ZoneSummaryReport,ZoneSummaryReportAdapter } from '@models/zone-summary-report';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { environment } from '@environments/environment';

@Injectable({
  providedIn: 'root'
})
export class ZoneService {
    constructor(private httpClient: HttpClient,
				private summary : ZoneSummaryReportAdapter,
				private climate: ZoneClimateReportAdapter) {
    }

    list(): Observable<ZoneSummaryReport[]> {
		return this.httpClient.get<any[]>(environment.apiEndpoint+'/zones').pipe(
			map(item => {
				let items = item as any[];
				let res: ZoneSummaryReport[]
				for ( let i of items ) {
					res.push(this.summary.adapt(i));
				}
				return res;
			}));
	}

	getZone(host: string, zone: string): Observable<ZoneClimateReport> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/host/'+host+'/zone/'+zone).pipe(
			map(item => {
				return this.climate.adapt(item)
			}));
	}
}
