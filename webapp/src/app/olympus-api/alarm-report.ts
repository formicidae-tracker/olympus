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

export class AlarmEvent {
  constructor(public start: Date, public end: Date | undefined) {}

  public time(): Date {
    return this.end == undefined ? this.start : this.end;
  }

  public on(): boolean {
    return this.end == undefined;
  }

  public duration(): number | undefined {
    if (this.end == undefined) {
      return undefined;
    }
    return this.end.getTime() - this.start.getTime();
  }

  static fromTimepoints(timepoints: any[]): AlarmEvent[] {
    let res: AlarmEvent[] = [];
    let lastOn: Date | undefined;
    for (const tp of timepoints) {
      if (lastOn != undefined) {
        if (tp.on == true) {
          continue;
        }
        res.push(new AlarmEvent(lastOn, new Date(tp.time)));
        lastOn = undefined;
      } else {
        if (tp.on == undefined || tp.on == false) {
          continue;
        }
        lastOn = new Date(tp.time);
      }
    }
    if (lastOn != undefined) {
      res.push(new AlarmEvent(lastOn, undefined));
    }

    return res;
  }
}

export class AlarmReport {
  public identification: string = '';
  public level: number = 0;
  public description: string = '';

  public events: AlarmEvent[] = [];

  static fromPlain(plain: any): AlarmReport {
    let res = new AlarmReport();
    res.identification = plain.identification || '';
    res.level = plain.level || 0;
    res.description = plain.description || '';
    res.events = AlarmEvent.fromTimepoints(plain.events || []);
    return res;
  }

  public on(): boolean {
    if (this.events.length == 0) {
      return false;
    }
    return this.events[this.events.length - 1].on();
  }

  public time(): Date {
    if (this.events.length == 0) {
      return new Date(0);
    }

    return this.events[this.events.length - 1].time();
  }

  public count(): number {
    return this.events.length;
  }

  static compareFunction(a: AlarmReport, b: AlarmReport): number {
    const res = AlarmReport._compareFunction(a, b);
    //console.log(a, b, res);
    return res;
  }

  static _compareFunction(a: AlarmReport, b: AlarmReport): number {
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
