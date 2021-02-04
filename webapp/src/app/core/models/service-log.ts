export class ServiceLog {
	constructor(public identifier: string = '',
				public on: boolean = false,
				public graceful: boolean = false,
				public time: Date = new Date()) {
	}

	static adapt(item: any): ServiceLog {
		return new ServiceLog(item.Identifier,item.On,item.Graceful,item.Time);
	}

	static adaptList(items: any[]): ServiceLog[] {
		if ( items == null ) {
			return [];
		}
		let res: ServiceLog[] = [];
		for ( let l of items ) {
			res.push(ServiceLog.adapt(l));
		}
		return res;
	}
}

export class ServiceLogReport {
	constructor(public climates: ServiceLog[][] = [],
				public trackings: ServiceLog[][] = []) {
	}
	static adapt(item: any): ServiceLogReport{
		let climates: ServiceLog[][] = [];
		if ( item.Climates != null ) {
			for ( let logs of item.Climates ) {
				climates.push(ServiceLog.adaptList(logs))
			}
		}
		let trackings: ServiceLog[][] = [];
		if ( item.Tracking != null ) {
			for ( let logs of item.Tracking ) {
				trackings.push(ServiceLog.adaptList(logs))
			}
		}
		return new ServiceLogReport(climates,trackings);
	}
}
