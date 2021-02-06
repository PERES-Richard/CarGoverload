import { Inject, Injectable } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';

@Injectable()
export class OffersService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka) { }

    startSearchingProcess(wishes: WishDTO[]) {
        this.kafkaClient.emit(`wish-requested`, wishes);
    }

}
