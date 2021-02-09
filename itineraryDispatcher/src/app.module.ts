import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { WishesModule } from './wishes/wishes.module';

@Module({
  imports: [WishesModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
