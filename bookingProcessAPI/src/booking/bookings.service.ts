import { Inject, Injectable } from '@nestjs/common';
import { ClientKafka } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import {Booking} from "../models/booking";

@Injectable()
export class BookingsService {
    constructor(@Inject("BOOKINGPROCESS_SERVICE") private readonly kafkaClient: ClientKafka) { }

    startBookValidation(wishes: WishDTO[]) {
        this.kafkaClient.emit(`book-validation`, wishes);
    }

    getWishesFromOfferID(offerID: string): WishDTO[] {
        // TODO search in redis
        return null;
    }

    payAndBookById(bookingID: string): boolean {
        // TODO
        let booking = this.getBookingFromID(bookingID)
        this.kafkaClient.emit(`book-register`, booking);
        return true
    }

    private getBookingFromID(bookingID: string): Booking {
        // TODO from redis
        return null
    }

    handleBookValidation(result) {
        // TODO
        // Get Offer
        // Compare offer
        // Create booking
    }
}
