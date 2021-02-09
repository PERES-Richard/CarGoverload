import { Controller, Logger } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { WishesService } from "./wishes.service";

@Controller('itinerary-dispatcher/wishes')
export class WishesController {
    constructor(private readonly wishesService: WishesService) {
    }
    @EventPattern("wish-requested")
    eventHandler(@Payload() data) {
        Logger.log(`Dispatching wishes...`)
        this.wishesService.dispatchWishes(data.value).then((_) => Logger.log("Dispatch done")).catch((err) => Logger.log(err));
    }
}
