import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { StreamInfo } from '@models/stream-info';
import { ZoneReport } from '@models/zone-report';

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})

export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string;
    hostName: string;
	zone: ZoneReport;
	update : Subscription;


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.zone = null;
	}

    ngOnInit() {
        this.zoneName = this.route.snapshot.paramMap.get('zoneName');
        this.hostName = this.route.snapshot.paramMap.get('hostName');
        this.title.setTitle('Olympus: '+this.hostName+'.'+this.zoneName)
		this.update = timer(0,5000).subscribe( (x) => {
			this.olympus.zoneReport(this.hostName,this.zoneName)
				.subscribe(
					(zone) => {
						this.zone = zone;
					},
					(error)  => {
						this.zone = null;
					},
					() => {

					});
		});
    }

	ngOnDestroy() {
		this.update.unsubscribe();
	}

}
