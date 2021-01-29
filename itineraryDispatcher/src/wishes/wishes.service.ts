import {Inject, Injectable} from '@nestjs/common';

@Injectable()
export class WishesService {
    constructor(@Inject("DISPATCHER_SERVICE") private readonly kafkaClient) {}

}
