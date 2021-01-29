import {Controller, Logger} from '@nestjs/common';
import {EventPattern, Payload} from "@nestjs/microservices";
import {WishesService} from "./wishes.service";

@Controller('itinerary-dispatcher/wishes')
export class WishesController {
    constructor(private readonly wishesService : WishesService) {
    }
    @EventPattern("wish-requested")
    eventHandler(@Payload() data){
        Logger.log(`The wish received is ${data.value}`);
    }
}
