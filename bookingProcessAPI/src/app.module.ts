import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { OffersModule } from './offers/offers.module';
import { RedisModule } from './redis/redis.module';

@Module({
  imports: [OffersModule, RedisModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule { }
