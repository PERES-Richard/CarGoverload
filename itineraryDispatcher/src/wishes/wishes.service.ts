import {Inject, Injectable} from '@nestjs/common';
import {ClientKafka} from "@nestjs/microservices";

@Injectable()
export class WishesService {
    constructor(@Inject("DISPATCHER_SERVICE") private readonly kafkaClient : ClientKafka) {}

}
