import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { WishResult } from 'src/models/wish_result';

@Injectable()
export class OffersService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka, @Inject("redis") private readonly redisClient) { }

    startSearchingProcess(wishes: WishDTO[]) {
        let wishRequest = { wishId: `${Date.now()}`, wishes: wishes };
        this.kafkaClient.emit(`wish-requested`, wishRequest);
        return wishRequest.wishId
    }

    async saveWishResult(result: WishResult) {
        Logger.log(`Saving result ${result.wishId}`);
        let isSaved = await this.redisClient.set(result.wishId, JSON.stringify(result.offerPossibilities));
        console.log(isSaved);
    }

}
