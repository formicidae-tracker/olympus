import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class TitleService {
  private _title = new Subject<string>();
  public title = this._title.asObservable();

  public setTitle(title: string): void {
    if (title.length == 0) {
      title = 'Olympus';
    } else {
      title = 'Olympus: ' + title;
    }
    this._title.next(title);
  }
}
