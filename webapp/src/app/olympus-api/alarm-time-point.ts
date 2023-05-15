import 'reflect-metadata';
import { Type, plainToClass } from 'class-transformer';

export class AlarmTimePoint {
  @Type(() => Date)
  public time: Date = new Date(0);

  public on: boolean = false;

  static fromPlain(plain: any): AlarmTimePoint {
    return plainToClass(AlarmTimePoint, plain, {
      exposeDefaultValues: true,
    });
  }
}
