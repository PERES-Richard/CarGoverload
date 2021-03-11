import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { OffersModule } from './offers/offers.module';
import { RedisModule } from './redis/redis.module';
import {BookingModule} from "./booking/booking.module";

@Module({
  imports: [OffersModule, RedisModule, BookingModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule { }
