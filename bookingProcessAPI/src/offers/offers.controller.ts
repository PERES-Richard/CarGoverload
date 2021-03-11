import { Body, Controller, Get, Logger, Param, Post } from '@nestjs/common';
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

    @Get('/:wishId')
    async offersResult(@Param('wishId') wishId: string) {
        return await this.offersService.getOffersSubjectOf(wishId);
    }

    @Post()
    getOffers(@Body() wishes: WishDTO[]): string {
        Logger.log(`Starting new search request`);
        const startedWish = this.offersService.startSearchingProcess(wishes);
        const res = {
            wishId: startedWish
        };
        return JSON.stringify(res);
    }
}
