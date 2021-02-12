import { Body, Controller, Get, Logger } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { OffersService } from "./offers.service";

@Controller('booking-process/offers')
export class OffersController {
    constructor(private readonly offersService: OffersService) {
    }
    @EventPattern("wish-result")
    wishResultHandler(@Payload() data) {
        Logger.log(`The wish result is ${data.value}`);
    }

    @Get()
    getOffers(@Body() wishes: WishDTO[]): string {
        Logger.log(`Starting new search request`);
        console.dir(wishes);
        let startedWish = this.offersService.startSearchingProcess(wishes);
        return `Search initiated, wishId generated : ${startedWish}`;
    }

    @Get('/payment')
    payOffer() {
        return true;
    }
}
