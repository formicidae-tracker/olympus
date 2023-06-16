export class NotificationSettings {
  public notifyOnWarning: boolean;
  public notifyNonGraceful: boolean;
  public subscribeToAll: boolean;
  public subscriptions: Set<string>;

  constructor({
    notifyNonGraceful = false,
    subscribeToAll = false,
    notifyOnWarning = false,
    subscriptions = new Set<string>([]),
  }: Partial<NotificationSettings> = {}) {
    this.notifyNonGraceful = notifyNonGraceful;
    this.subscribeToAll = subscribeToAll;
    this.notifyOnWarning = notifyOnWarning;
    this.subscriptions = subscriptions;
  }

  public toPlain(): any {
    let plain: any = Object.assign({}, this);
    plain.subscriptions = Array.from(this.subscriptions);
    return plain;
  }

  public serialize(): string {
    return JSON.stringify(this.toPlain());
  }

  static fromPlain(plain: any): NotificationSettings {
    // make a deep copy to preserve the plain object
    plain = Object.assign({}, plain);
    const subscriptions = plain.subscriptions;
    if (subscriptions != undefined) {
      plain.subscriptions = new Set<string>(subscriptions);
    }
    return new NotificationSettings(plain as Partial<NotificationSettings>);
  }

  static deserialize(jsondata: string) {
    return NotificationSettings.fromPlain(JSON.parse(jsondata));
  }

  public hasSubscription(zone: string): boolean {
    return this.subscribeToAll || this.subscriptions.has(zone);
  }

  public subscribeTo(zone: string): boolean {
    if (this.subscribeToAll || this.subscriptions.has(zone)) {
      return false;
    }
    this.subscriptions.add(zone);
    return true;
  }

  public unsubscribeFrom(zone: string): boolean {
    if (this.subscribeToAll || this.subscriptions.has(zone) == false) {
      return false;
    }
    this.subscriptions.delete(zone);
    return true;
  }

  public needPushSubscription(): boolean {
    return this.subscribeToAll || this.subscriptions.size > 0;
  }
}
