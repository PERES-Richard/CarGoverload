import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishPaymentDTO } from "../models/wish_payment_dto";

@Injectable()
export class BookingsService {
    validationsResult: Map<string, string> = new Map();

    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka) { }

    payAndBookById(wishPayment: WishPaymentDTO): boolean {
        this.kafkaClient.emit(`book-validation`, wishPayment);
        return true
    }

    handleBookConfirmation(result: any) {
        Logger.log("Booking confirmation received : " + result);
        console.dir(result)
        this.validationsResult.set(result.wishId, result.result);
    }
}
