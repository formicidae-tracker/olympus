import { AlarmTimePoint } from './alarm-time-point';

function compareTimeDescending(a: Date, b: Date): number {
  if (a == b) {
    return 0;
  }
  return a > b ? -1 : 1;
}

function compareLevelDescending(a: number, b: number): number {
  if (a == b) {
    return 0;
  }
  return a > b ? -1 : 1;
}

function compareOnDescending(a: boolean, b: boolean): number {
  if (a == b) {
    return 0;
  }
  return a ? -1 : 1;
}

export class AlarmReport {
  public identification: string = '';
  public level: number = 0;
  public description: string = '';

  public events: AlarmTimePoint[] = [];

  static fromPlain(plain: any): AlarmReport {
    let res = new AlarmReport();
    res.identification = plain.identification || '';
    res.level = plain.level || 0;
    res.description = plain.description || '';
    for (const e of plain.events || []) {
      res.events.push(AlarmTimePoint.fromPlain(e));
    }
    return res;
  }

  public on(): boolean {
    if (this.events.length == 0) {
      return false;
    }
    return this.events[this.events.length - 1].on;
  }

  public time(): Date {
    if (this.events.length == 0) {
      return new Date(0);
    }
    return this.events[this.events.length - 1].time;
  }

  public since(t: Date): string {
    let ellapsed: number = t.getTime() - this.time().getTime();
    if (ellapsed <= 10000) {
      return 'now';
    }

    for (const t of thresholds) {
      ellapsed = Math.round(ellapsed / t.divider);
      if (t.threshold == undefined || ellapsed < t.threshold) {
        return ellapsed + t.units;
      }
    }

    return 'never';
  }

  public count(): number {
    let res: number = 0;
    for (const p of this.events) {
      if (p.on) {
        res += 1;
      }
    }
    return res;
  }

  static compareFunction(a: AlarmReport, b: AlarmReport): number {
    let compareOn = compareOnDescending(a.on(), b.on());
    if (compareOn != 0) {
      return compareOn;
    }
    let compareLevel = compareLevelDescending(a.level, b.level);
    if (compareLevel != 0) {
      return compareLevel;
    }
    return compareTimeDescending(a.time(), b.time());
  }
}

class Threshold {
  constructor(
    public units: string,
    public divider: number,
    public threshold?: number
  ) {}
}

const thresholds: Threshold[] = [
  new Threshold('s', 1000, 60),
  new Threshold('m', 60, 60),
  new Threshold('h', 60, 24),
  new Threshold('d', 24, 7),
  new Threshold('w', 7, undefined),
];
