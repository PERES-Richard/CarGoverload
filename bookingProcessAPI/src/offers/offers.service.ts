import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { BehaviorSubject, Observable } from 'rxjs';
import { WishDTO } from 'src/models/wish_dto';
import { OfferPossibility, WishResult } from 'src/models/wish_result';

@Injectable()
export class OffersService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka, @Inject("redis") private readonly redisClient) { }

    offersResults: Map<string, BehaviorSubject<OfferPossibility[]>> = new Map();

    startSearchingProcess(wishes: WishDTO[]) {
        const wishRequest = { wishId: `${Date.now()}`, wishes: wishes };
        this.kafkaClient.emit(`wish-requested`, wishRequest);
        this.offersResults.set(wishRequest.wishId, new BehaviorSubject(null));
        return wishRequest.wishId
    }

    async saveWishResult(result: WishResult) {
        Logger.log(`Saving result ${result.wishId}`);
        this.redisClient.set(result.wishId, JSON.stringify(result.offerPossibilities))
            .then((done) => {
                console.log(done);
                this.offersResults.get(result.wishId).next(result.offerPossibilities);

            })
            .catch((err) => Logger.log(err));

    }

    getOffersSubjectOf(wishId: string): Observable<OfferPossibility[]> {
        return this.offersResults.get(wishId).asObservable();
    }

}
