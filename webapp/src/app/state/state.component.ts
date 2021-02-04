import { Component, OnInit, Input } from '@angular/core';
import { State }  from '@models/state';


@Component({
    selector: 'app-state',
    templateUrl: './state.component.html',
    styleUrls: ['./state.component.css']
})
export class StateComponent implements OnInit {
    @Input() current: State;
	@Input() end: State;
	@Input() currentTemperature: number;
	@Input() currentHumidity: number;
	@Input() displayCurrent: boolean;


	static displayValue(v: number) :string {
		if ( v == null ||isNaN(v) == true ) {
			return 'N.A.';
		}
		return (Math.round(100*v)/100).toString();
	}

	displayField(field): string {
		if ( this.current == null ) {
			return 'N.A.';
		}
		if (this.end == null) {
			return StateComponent.displayValue(this.current[field]);
		}
		return StateComponent.displayValue(this.current[field]) + ' ↦ ' + StateComponent.displayValue(this.end[field]);
	}

    constructor() {
		this.current = null;
		this.end = null;
		this.currentTemperature = NaN;
		this.currentHumidity = NaN;
		this.displayCurrent = false;
	}

    ngOnInit() {
    }

}
