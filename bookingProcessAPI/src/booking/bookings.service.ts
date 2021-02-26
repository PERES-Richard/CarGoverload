import { Inject, Injectable } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import {WishPaymentDTO} from "../models/wish_payment_dto";

@Injectable()
export class BookingsService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka) { }

    payAndBookById(wishPayment: WishPaymentDTO): boolean {
        this.kafkaClient.emit(`book-validation`, wishPayment);
        return true
    }

    handleBookValidation(result) {
        // TODO
        // Get Offer
        // Compare offer
        // Create booking
    }
}
