import { Injectable } from '@angular/core';
import { Adapter } from './adapter';

import { ZoneClimateStatus } from './zone-climate-status';


export class ZoneSummary {
	constructor(public Host: string,
				public ZoneName: string,
				public StreamURL: string  = '',
				public Status: ZoneClimateStatus = null ) {
	}
}
