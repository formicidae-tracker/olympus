export type EventStatus = 'running' | 'done' | 'ungraceful';

export class Event {
  constructor(
    public start: Date = new Date(0),
    public end?: Date,
    public graceful?: boolean
  ) {}

  public time(): Date {
    return this.end || this.start;
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

  public status(): EventStatus {
    if (this.end == undefined) {
      return 'running';
    }
    if (this.graceful == undefined || this.graceful == true) {
      return 'done';
    }
    return 'ungraceful';
  }

  static fromPlain(plain: any): Event {
    let end = undefined;
    if (plain.end != undefined) {
      end = new Date(plain.end);
    }
    return new Event(new Date(plain.start || 0), end, plain.graceful);
  }
}
