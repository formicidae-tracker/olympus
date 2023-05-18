export class UserSettings {
  public darkMode: boolean;
  public alarmSubscriptions: Set<string>;

  constructor({
    darkMode = false,
    alarmSubscriptions = new Set<string>([]),
  }: Partial<UserSettings> = {}) {
    this.darkMode = darkMode;
    this.alarmSubscriptions = alarmSubscriptions;
  }

  public serialize(): string {
    let plain = {
      darkMode: this.darkMode,
      alarmSubscriptions: Array.from(this.alarmSubscriptions),
    };
    return JSON.stringify(plain);
  }

  static deserialize(jsondata: string) {
    let plain = JSON.parse(jsondata) as Partial<UserSettings>;
    if (plain.alarmSubscriptions != undefined) {
      plain.alarmSubscriptions = new Set<string>(plain.alarmSubscriptions);
    }
    return new UserSettings(plain);
  }

  public hasSubscriptionToAlarmFrom(zone: string): boolean {
    return this.alarmSubscriptions.has(zone);
  }

  public subscribeToAlarmFrom(zone: string): boolean {
    if (this.alarmSubscriptions.has(zone)) {
      return false;
    }
    this.alarmSubscriptions.add(zone);
    return true;
  }

  public unsubscribeFromAlarmFrom(zone: string): boolean {
    if (this.alarmSubscriptions.has(zone) == false) {
      return false;
    }
    this.alarmSubscriptions.delete(zone);
    return true;
  }
}
