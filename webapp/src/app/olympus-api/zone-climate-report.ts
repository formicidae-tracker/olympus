import { ClimateState } from './climate-state';
import { Bounds } from './bounds';

export class ZoneClimateReport {
  since: Date = new Date(0);
  temperature?: number;
  humidity?: number;

  temperature_bounds?: Bounds;
  humidity_bounds?: Bounds;

  current?: ClimateState;
  current_end?: ClimateState;

  next?: ClimateState;
  next_end?: ClimateState;
  next_time?: Date;

  static fromPlain(plain: any): ZoneClimateReport | undefined {
    if (plain == undefined) {
      return undefined;
    }
    let ret = new ZoneClimateReport();
    ret.temperature = plain.temperature;
    ret.humidity = plain.humidity;
    ret.since = new Date(plain.since || 0);
    if (plain.temperature_bounds != undefined) {
      ret.temperature_bounds = Bounds.fromPlain(plain.temperature_bounds);
    }
    if (plain.humidity_bounds != undefined) {
      ret.humidity_bounds = Bounds.fromPlain(plain.humidity_bounds);
    }

    if (plain.current != undefined) {
      ret.current = ClimateState.fromPlain(plain.current);
    }
    if (plain.current_end != undefined) {
      ret.current_end = ClimateState.fromPlain(plain.current_end);
    }

    if (plain.next != undefined) {
      ret.next = ClimateState.fromPlain(plain.next);
    }
    if (plain.next_end != undefined) {
      ret.next_end = ClimateState.fromPlain(plain.next_end);
    }
    if (plain.next_time != undefined) {
      ret.next_time = new Date(plain.next_time);
    }

    return ret;
  }
}
