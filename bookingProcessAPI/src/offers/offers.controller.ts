import { Body, Controller, Logger, Param, Post, Sse } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { OffersService } from "./offers.service";

@Controller('booking-process/offers')
export class OffersController {
    constructor(private readonly offersService: OffersService) {
    }

    @EventPattern("wish-result")
    wishResultHandler(@Payload() data) {
        Logger.log(`The wish result is ${data.value.wishId} has been received, saving...`);
        this.offersService.saveWishResult(data.value);
    }

    @Sse('/:wishId')
    offersResult(@Param('wishId') wishId: string) {
        Logger.log(`New sse client on ${wishId}`);
        return this.offersService.getOffersSubjectOf(wishId);
    }

    @Post()
    getOffers(@Body() wishes: WishDTO[]): string {
        Logger.log(`Starting new search request`);
        console.dir(wishes);
        const startedWish = this.offersService.startSearchingProcess(wishes);
        const res = {
            wishId: startedWish
        };
        return JSON.stringify(res);
    }
}
