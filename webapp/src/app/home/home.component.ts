import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { ZoneService } from '@services/zone';
import { ZoneSummaryReport }  from '@models/zone-summary-report';
import { Subscription,timer } from 'rxjs';
@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit,OnDestroy {
    zones: ZoneSummaryReport[];
	update: Subscription;


    constructor(private zs : ZoneService, private title: Title) {
		this.zones = [];
	}

    ngOnInit() {
		this.title.setTitle('Olympus: Home')

		this.update = timer(0,20000).subscribe(x => {
			this.zs.list().subscribe( (list) => {
				this.zones = list;
			});
		})
    }

	ngOnDestroy() {
		this.update.unsubscribe()
	}

}
