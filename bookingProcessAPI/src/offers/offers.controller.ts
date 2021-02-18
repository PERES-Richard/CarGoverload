import {Body, Controller, Logger, Post} from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { OffersService } from "./offers.service";

@Controller('booking-process/offers')
export class OffersController {
    constructor(private readonly offersService: OffersService) {
    }
    @EventPattern("wish-result")
    wishResultHandler(@Payload() data) {
        // TODO stock offers result of this wish in redis
        // TODO print wish id also
        Logger.log(`The wish result is ${data.value}`);
    }

    @Post()
    getOffers(@Body() wishes: WishDTO[]): string {
        Logger.log(`Starting new search request`);
        console.dir(wishes);
        const startedWish = this.offersService.startSearchingProcess(wishes);
        return `Search initiated, wishId generated : ${startedWish}`;
    }
}
