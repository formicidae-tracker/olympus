import { Injectable } from '@angular/core';
import { ZoneSummaryReport } from '@models/zone-summary-report';
import { ClimateTimeSeries } from '@models/climate-time-series';
import { StreamInfo } from '@models/stream-info';
import { Observable,throwError,of } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';
import { environment } from '@environments/environment';
import { State } from '@models/state';
import { Bounds } from '@models/bounds';
import { ZoneReport } from '@models/zone-report';
import { ZoneClimateReport } from '@models/zone-climate-report';



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


export class MockOlympusService {
	constructor() {
	}

	private staticData = {
		somehost: {
			box: {
				climate: new ZoneClimateReport(
					21.0,
					61.0,
					new Bounds(15,30),
					new Bounds(40,80),
					2,
					1,
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
				alarms: null,
				streamInfo: new StreamInfo('/olympus/hls/somehost.m3u8',
										   '/olympus/somehost.png'),
			},
			tunnel: {
				climate: new ZoneClimateReport(
					-1000.0,
						-1000.0,
					new Bounds(0,100),
					new Bounds(0,100),
					0,
					0,
					0,
					new State("always-on",-1000,-1000,100,100,-1000),
					null,
					null,
					null,
					null
				),
				alarms: null,
				timeSeries: new ClimateTimeSeries(
					null,
					null,
					null
				),
				streamInfo: null
			},
		},
		notracking: {
			box: {
				climate: new ZoneClimateReport(
					21.0,
					61.0,
					new Bounds(15,30),
					new Bounds(40,80),
					2,
					1,
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
				streamInfo: null,
				alarms: null,
			},
		},
		onlytracking: {
			box: {
				climate: null,
				alarms: null,
				timeSeries: null,
				streamInfo: new StreamInfo('/olympus/hls/onlytracking.m3u8','/olympus/onlytracking.png'),
			}
		},
	}

	zoneReportStatic(host: string,zone: string): ZoneReport {
		if ( this.staticData[host] == null
			|| this.staticData[host][zone] == null) {
			return null;
		}
		let z = this.staticData[host][zone];
		return new ZoneReport(host,zone,z.climate,z.streamInfo,z.alarms);
	}

	zoneReport(host: string,zone: string): Observable<ZoneReport> {
		let res = this.zoneReportStatic(host,zone);
		if ( res == null ) {
			return throwError('mock olympus: zoneReport: unknown zone '+host+'/zone/'+zone);
		}
		return of(res);
	}


	zoneSummariesStatic(): ZoneSummaryReport[] {
		return [
			new ZoneSummaryReport('notracking',
								  'box',
								  new StreamInfo(),
								  this.staticData.notracking.box.climate),
			new ZoneSummaryReport('onlytracking',
								  'box',
								  new StreamInfo('/olympus/hls/onlytracking.m3u8','/olympus/onlytracking.png'),
								  null),
			new ZoneSummaryReport('somehost',
								  'box',
								  new StreamInfo('/olympus/hls/somehost.m3u8','/olympus/somehost.png'),
								  this.staticData.somehost.box.climate),
			new ZoneSummaryReport('somehost',
								  'tunnel',
								  new StreamInfo(),
								  this.staticData.somehost.tunnel.climate),

		];
	}
	zoneSummaries(): Observable<ZoneSummaryReport[]> {
		return of(this.zoneSummariesStatic());
	}

	climateTimeSeriesStatic(host: string, zone: string, window: string): ClimateTimeSeries {
		if ( this.staticData[host] == null
			|| this.staticData[host][zone] == null
			|| this.staticData[host][zone].timeSeries == null ) {
			return null;
		}
		return this.staticData[host][zone].timeSeries;
	}

	climateTimeSeries(host: string, zone: string, window: string): Observable<ClimateTimeSeries> {
		let res = this.climateTimeSeriesStatic(host,zone,window);
		if ( res == null ) {
			return throwError('mock olympus: climateSeries: unknown zone '+host+'/zone/'+zone);
		}
		return of(res);
	}

}
