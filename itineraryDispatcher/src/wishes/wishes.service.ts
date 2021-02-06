import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';

@Injectable()
export class WishesService {
    constructor(@Inject("DISPATCHER_SERVICE") private readonly kafkaClient: ClientKafka) { }

    async dispatchWishes(wishes: WishDTO[]) {
        const requestTimestamp = Date.now();
        for (let indexOfWishes = 0; indexOfWishes < wishes.length; indexOfWishes++) {
            const wish = wishes[indexOfWishes];
            for (let indexOfCars = 0; indexOfCars < wish.numberOfCars; indexOfCars++) {
                let formatedWish = { id: `${requestTimestamp}_${indexOfWishes}-${wish.carType.charAt(0).toUpperCase()}${indexOfCars}`, carType: wish.carType, departureNode: wish.departureNode, arrivalNode: wish.arrivalNode };
                await this.kafkaClient.emit(`new-search`, formatedWish).toPromise();

            }

        }
    }

}
