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

  public serialize(): string {
    let plain: any = Object.assign({}, this);
    plain.subscriptions = Array.from(this.subscriptions);
    return JSON.stringify(plain);
  }

  static deserialize(jsondata: string) {
    let plain = JSON.parse(jsondata);
    const subscriptions = plain.subscriptions;
    if (subscriptions != undefined) {
      plain.subscriptions = new Set<string>(subscriptions);
    }
    return new NotificationSettings(plain as Partial<NotificationSettings>);
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
}
