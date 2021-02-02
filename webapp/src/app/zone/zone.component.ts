import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { ZoneReport } from '@models/zone-report';
import { switchMap } from 'rxjs/operators';

export enum State {
	Loading = 1,
	Loaded = 2,
	Unavailable = 3,
};

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})


export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string;
    hostName: string;
	zone: ZoneReport;
	state: State;
	update: Subscription;


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.state = State.Loading;
		this.zone = null;
		this.update = null;
	}

    ngOnInit() {
		this.route.paramMap.pipe(
			.switchMap(params => {
				let hostName = params.get('hostName');
				let zoneName = params.get('zoneName');
				this.title.setTitle('Olympus: '+hostName+'.'+zoneName)
				this.update = timer(0,5000).subscribe( () => {
					this.olympus.zoneReport(hostName,zoneName)
						.subscribe(
							(zone) => {
								this.zone = zone;
								this.state = State.Loaded;
							},
							()  => {
								this.zone = null;
								this.state = State.Unavailable;
							},
							() => {

							});
				});
			}),
		);
    }

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}

}
