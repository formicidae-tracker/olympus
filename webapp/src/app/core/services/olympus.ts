import { Injectable } from '@angular/core';
import { ZoneClimateReport } from '@models/zone-climate-report';
import { ZoneSummaryReport } from '@models/zone-summary-report';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { StreamInfo } from '@models/stream-info';
import { Observable,throwError,of } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { environment } from '@environments/environment';
import { ZoneClimateStatus } from '@models/zone-climate-status';
import { State } from '@models/state';
import { Bounds } from '@models/bounds';



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

	private staticData = {
		somehost: {
			box: {
				climate: new ZoneClimateReport(
					new ZoneClimateStatus(21.0,61.0,new Bounds(15,30),new Bounds(40,80),2,1),
					0,
					new State("day",60.0,21.5,100,100,100),
					null,
					new State("day to night",60.0,21.5,100,100,100),
					new State("night",60.0,19.0,100,0,0),
					new Date()
				),
				timeSeries: new ClimateTimeSeries(
					[{X:0,Y:61.2},{X:1,Y:61.0}],
					[{X:0,Y:21.4},{X:1,Y:21.5}],
					null
				),
				streamInfo: new StreamInfo('https://olympus.com/olympus/hls/somehost.m3u8'),
			},
			tunnel: {
				climate: new ZoneClimateReport(
					new ZoneClimateStatus(-1000.0,-1000.0,new Bounds(0,100),new Bounds(0,100),0,0),
					0,
					new State("always-on",-1000,-1000,100,100,-1000),
					null,
					null,
					null,
					null
				),
				timeSeries: new ClimateTimeSeries(
					null,
					null,
					null
				),
			},
		},
		notracking: {
			box: {
				climate: new ZoneClimateReport(
					new ZoneClimateStatus(21.0,61.0,new Bounds(15,30),new Bounds(40,80),2,1),
					0,
					new State("day",60.0,21.5,100,100,100),
					null,
					new State("day to night",60.0,21.5,100,100,100),
					new State("night",60.0,19.0,100,0,0),
					new Date()
				),
				timeSeries: new ClimateTimeSeries(
					[{X:0,Y:61.2},{X:1,Y:61.0}],
					[{X:0,Y:21.4},{X:1,Y:21.5}],
					null
				),
				streamInfo: new StreamInfo(''),
			},
		},
		onlytracking: {
			climate: {}
		},

	}

	zoneClimate(host: string,zone: string): Observable<ZoneClimateReport> {
		if ( this.staticData[host] == null
			|| this.staticData[host][zone] == null) {
			return throwError('olympus: unknown zone '+host+'/zone/'+zone);
		}
		return of(this.staticData[host][name]['climate'])
	}

	zoneSummaries(): Observable<ZoneSummaryReport[]> {
		return of([
			new ZoneSummaryReport('notracking',
								  'box',
								  new StreamInfo(''),
								  this.staticData.notracking.box.climate.ClimateStatus),
			new ZoneSummaryReport('onlytracking',
								  'box',
								  new StreamInfo('https://olympus.com/olympus/hls/onlytracking.m3u8'),
								  null),
			new ZoneSummaryReport('somehost',
								  'box',
								  new StreamInfo('https://olympus.com/olympus/hls/somehost.m3u8'),
								  this.staticData.somehost.box.climate.ClimateStatus),
			new ZoneSummaryReport('somehost',
								  'tunnel',
								  new StreamInfo(''),
								  this.staticData.somehost.tunnel.climate.ClimateStatus),

		]);
	}

	streamURL(host: string): Observable<StreamInfo> {
		if ( this.staticData[host] == null || this.staticData[host].streamInfo.streamURL.length == 0 ) {
			return throwError('no stream info for '+host);
		}
		return of(this.staticData[host].streamInfo);
	}

	climateTimeSeries(host: string, zone: string, window: string): Observable<ClimateTimeSeries> {
		return ;
	}
}
