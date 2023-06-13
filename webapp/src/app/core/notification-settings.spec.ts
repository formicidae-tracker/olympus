import { cases } from 'jasmine-parameterized';
import { NotificationSettings } from './notification-settings';

describe('NotificationSettings', () => {
  it('should becreated with default values', () => {
    let e = new NotificationSettings();
    expect(e).toBeTruthy();
    expect(e.subscriptions).toEqual(new Set<string>([]));
  });

  it('should be created with non-default values', () => {
    let e = new NotificationSettings({
      subscriptions: new Set<string>(['some.zone']),
    });
    expect(e).toBeTruthy();
    expect(e.subscriptions).toEqual(new Set<string>(['some.zone']));
  });

  cases([
    [
      {},
      '{"notifyNonGraceful":false,"subscribeToAll":false,"notifyOnWarning":false,"subscriptions":[]}',
    ],
    [
      {
        notifyNonGraceful: true,
        subscribeToAll: true,
        notifyOnWarning: true,
        subscriptions: new Set<string>(['some.zone']),
      },
      '{"notifyNonGraceful":true,"subscribeToAll":true,"notifyOnWarning":true,"subscriptions":["some.zone"]}',
    ],
  ]).it('should serialize to JSON', ([args, jsondata]) => {
    expect(new NotificationSettings(args).serialize()).toEqual(jsondata);
  });

  cases([
    [{}, '{}'],
    [{}, '{"notifyNonGraceful":false,"subscriptions":[]}'],
    [
      {
        notifyNonGraceful: true,
        subscriptions: new Set<string>(['some.zone']),
      },
      '{"notifyNonGraceful":true,"subscriptions":["some.zone"]}',
    ],
    [{}, '{"oldKey":true}'],
  ]).it('should deserialize from JSON', ([args, jsondata]) => {
    expect(NotificationSettings.deserialize(jsondata)).toEqual(
      new NotificationSettings(args)
    );
  });

  it('should handle simple alarm subscription', () => {
    let s = new NotificationSettings();
    expect(s.hasSubscription('foo')).toBeFalse();
    expect(s.subscribeTo('foo')).toBeTrue();
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.subscribeTo('foo')).toBeFalse();
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.unsubscribeFrom('foo')).toBeTrue();
    expect(s.hasSubscription('foo')).toBeFalse();
    expect(s.unsubscribeFrom('foo')).toBeFalse();
  });

  it('should handle subscribeToAll', () => {
    let s = new NotificationSettings();
    expect(s.hasSubscription('foo')).toBeFalse();
    s.subscribeToAll = true;
    expect(s.hasSubscription('foo')).toBeTrue();
    expect(s.subscribeTo('foo')).toBeFalse();
    expect(s.unsubscribeFrom('foo')).toBeFalse();
    expect(s.subscriptions).toHaveSize(0);
    s.subscribeToAll = false;
    expect(s.hasSubscription('foo')).toBeFalse();
  });

  cases([
    [{}, false],
    [{ notifyOnWarning: true }, false],
    [{ subscribeToAll: true }, true],
    [{ subscriptions: new Set<string>(['foo']) }, true],
  ]).it('should inform if it needs susbscription', ([args, expected]) => {
    const s = new NotificationSettings(args);
    expect(s.needPushSubscription()).toEqual(expected);
  });
});
