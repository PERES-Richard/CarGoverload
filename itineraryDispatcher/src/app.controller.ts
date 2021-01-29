import {Controller, Get, Logger} from '@nestjs/common';
import { AppService } from './app.service';
import {EventPattern, Payload} from "@nestjs/microservices";

@Controller("itinerary-dispatcher")
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  getHello(): string {
    return this.appService.getHello();
  }

  @EventPattern("wish-requested")
  eventHandler(@Payload() data){
    Logger.log(`The message received is ${data.value}`);
  }
}
