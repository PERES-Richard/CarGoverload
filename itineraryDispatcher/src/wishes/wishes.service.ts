import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { SearchEvent } from 'src/models/search_event';
import { WishDTO } from 'src/models/wish_dto';
import { WishEvent } from 'src/models/wish_event';

@Injectable()
export class WishesService {
    constructor(@Inject("DISPATCHER_SERVICE") private readonly kafkaClient: ClientKafka) { }

    async dispatchWishes(wishes: WishDTO[]) {
        const requestTimestamp = Date.now();
        let wishEvent: WishEvent = { wishId: `${requestTimestamp}`, searchIds: [] };
        for (let indexOfWishes = 0; indexOfWishes < wishes.length; indexOfWishes++) {
            const wish = wishes[indexOfWishes];
            for (let indexOfCars = 0; indexOfCars < wish.numberOfCars; indexOfCars++) {
                let searchEvent: SearchEvent = { searchId: `${requestTimestamp}_${indexOfWishes}-${wish.carType.charAt(0).toUpperCase()}${indexOfCars}`, carType: wish.carType, departureNode: wish.departureNode, arrivalNode: wish.arrivalNode, dateDeparture: wish.dateDeparture };
                wishEvent.searchIds.push(searchEvent.searchId);
                await this.kafkaClient.emit(`new-search`, searchEvent).toPromise();
            }

        }
        this.kafkaClient.emit(`new-wish`, wishEvent).toPromise()
            .then((_) => Logger.log(`Wishes of request ${wishEvent.wishId} sent on new-wish topic`))
            .catch((err) => Logger.log(`Error on wishes request ${wishEvent.wishId} : ${err}`));
    }

}
