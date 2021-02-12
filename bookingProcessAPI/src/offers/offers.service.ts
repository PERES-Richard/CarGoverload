import { Inject, Injectable } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { RedisClient } from 'redis';
import { WishDTO } from 'src/models/wish_dto';

@Injectable()
export class OffersService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka, @Inject("redis") private readonly redisClient: RedisClient) { }

    startSearchingProcess(wishes: WishDTO[]) {
        let wishRequest = { wishId: `${Date.now()}`, wishes: wishes };
        this.kafkaClient.emit(`wish-requested`, wishRequest);
        return wishRequest.wishId
    }

}
