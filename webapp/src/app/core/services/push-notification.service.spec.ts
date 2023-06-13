import { TestBed } from '@angular/core/testing';

import { PushNotificationService } from './push-notification.service';
import { Observable, of } from 'rxjs';
import { SwPush } from '@angular/service-worker';
import { HttpClientModule } from '@angular/common/http';

export class FakeSwPush {
  public subscription: Observable<PushSubscription | null> = of(null);
}

describe('PushNotificationService', () => {
  let service: PushNotificationService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientModule],
      providers: [{ provide: SwPush, useClass: FakeSwPush }],
    });
    service = TestBed.inject(PushNotificationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
