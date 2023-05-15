import { Expose, plainToClass } from 'class-transformer';

export class Bounds {
  @Expose() public minimum?: number;
  @Expose() public maximum?: number;

  static fromPlain(plain: any): Bounds {
    return plainToClass(Bounds, plain, {
      exposeDefaultValues: true,
    });
  }
}
