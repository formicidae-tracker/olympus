import { Injectable } from '@angular/core';
import { ZoneClimateReport } from '@models/zone-climate-report';
import { ZoneSummaryReport } from '@models/zone-summary-report';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { StreamInfo } from '@models/stream-info';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { environment } from '@environments/environment';



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
				let res: ZoneSummaryReport[]
				for ( let i of items ) {
					res.push(ZoneSummaryReport.adapt(i));
				}
				return res;
			}));
	}

	zoneClimate(host: string, zone: string): Observable<ZoneClimateReport> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/host/'+host+'/zone/'+zone).pipe(
			map(item => {
				return ZoneClimateReport.adapt(item)
			}));
	}

	streamURL(host: string): Observable<StreamInfo> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/tracking/host/'+host).pipe(
			map(item => {
				return new StreamInfo(item.StreamURL);
			}),
		);
	}

	climateTimeSeries(host: string, zone: string, window: string): Observable<ClimateTimeSeries> {
		return this.httpClient.get<any>(environment.apiEndpoint+'/host/'+host+'/zone/'+zone+'/climate?window='+window).pipe(
			map(item => { return ClimateTimeSeries.adapt(item); })
		);
	}

}


export class MockOlympusService {
	constructor() {
	}

	zoneClimate(host: string,zone: string): Observable<ZoneClimateReport> {
		return null;
	}

	zoneSummaries(): Observable<ZoneSummaryReport> {
		return null;
	}

	streamURL(host: string): Observable<StreamInfo> {
		return null;
	}

	climateTimeSeries(host: string, zone: string, window: string): Observable<ClimateTimeSeries> {
		return null;
	}
}
