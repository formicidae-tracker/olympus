import { Component, OnInit, Input } from '@angular/core';
import { State }  from '@models/state';


@Component({
    selector: 'app-state',
    templateUrl: './state.component.html',
    styleUrls: ['./state.component.css']
})
export class StateComponent implements OnInit {
    @Input() stateA: State;
	@Input() stateB: State;
	@Input() currentTemperature: number;
	@Input() currentHumidity: number;
	@Input() displayCurrent: boolean;


	displayValue(v :number) :string {
		if ( isNaN(v) == true ) {
			return 'N.A.';
		}
		return (Math.round(100*v)/100).toString();
	}

    constructor() {
		this.stateA = new State();
		this.stateB = new State();
		this.currentTemperature = NaN;
		this.currentHumidity = NaN;
		this.displayCurrent = false;
	}

    ngOnInit() {
    }

}
