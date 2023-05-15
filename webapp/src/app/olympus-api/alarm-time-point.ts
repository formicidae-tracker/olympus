export class AlarmTimePoint {
  public time: Date = new Date(0);

  public on: boolean = false;

  static fromPlain(plain: any): AlarmTimePoint {
    let res = new AlarmTimePoint();
    res.time = new Date(plain.time || 0);
    res.on = plain.on || false;
    return res;
  }
}
