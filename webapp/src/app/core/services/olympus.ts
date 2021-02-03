import { Injectable } from '@angular/core';
import { ZoneSummaryReport } from '@models/zone-summary-report';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { environment } from '@environments/environment';
import { ZoneReport } from '@models/zone-report';



@Injectable({
  providedIn: 'root'
})

export class OlympusService {
    constructor(private httpClient: HttpClient) {
    }

    zoneSummaries(): Observable<ZoneSummaryReport[]> {
		return this.httpClient.get<any[]>(environment.apiEndpoint+'/zones').pipe(
			map(item => {
				let items = item as any[];
				let res: ZoneSummaryReport[] = [];
				for ( let i of items ) {
					res.push(ZoneSummaryReport.adapt(i));
				}
				return res;
			}));
	}

	zoneReport(host: string, zone: string): Observable<ZoneReport> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/host/'+host+'/zone/'+zone).pipe(
			map(item => {
				return ZoneReport.adapt(item)
			}));
	}

	climateTimeSeries(host: string, zone: string, window: string): Observable<ClimateTimeSeries> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/host/'+host+'/zone/'+zone+'/climate?window='+window).pipe(
			map(item => { return ClimateTimeSeries.adapt(item); })
		);
	}

}
