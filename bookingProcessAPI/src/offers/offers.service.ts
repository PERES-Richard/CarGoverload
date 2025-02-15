import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { WishResult } from 'src/models/wish_result';

@Injectable()
export class OffersService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka,
                @Inject("redis") private readonly redisClient: any) { }

   // offersResults: Map<string, BehaviorSubject<WishResult>> = new Map();
    offersResults: Map<string, WishResult> = new Map();

    startSearchingProcess(wishes: WishDTO[]) {
        const wishRequest = { wishId: `${Date.now()}`, wishes: wishes };
        this.kafkaClient.emit(`wish-requested`, wishRequest);
        // this.offersResults.set(wishRequest.wishId, new BehaviorSubject(null));
        this.offersResults.set(wishRequest.wishId, null);
        return wishRequest.wishId
    }

    async saveWishResult(result: WishResult) {
        Logger.log(`Saving result ${result.wishId}`);
        this.redisClient.set(result.wishId, JSON.stringify(result))
            .then((done) => {
                console.log(done);
                //this.offersResults.get(result.wishId).next(result);
            })
            .catch((err) => Logger.log(err));
    }

    async getOffersSubjectOf(wishId: string) {
        return await this.redisClient.get(wishId);
        //return this.offersResults.get(wishId);
    }

}
