import { Body, Controller, Post, Logger } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { BookingsService } from "./bookings.service";
import {WishPaymentDTO} from "../models/wish_payment_dto";

@Controller('booking-process/booking/')
export class BookingsController {
    constructor(private readonly bookingsService: BookingsService) {
    }
    
    @EventPattern("book-validation-result")
    bookValidationResultHandler(@Payload() data) {
        Logger.log(`Book validation result received. Processing validity...`);
        // TODO Check if offer selected is in the Offer's List
        let result
        this.bookingsService.handleBookValidation(result)
    }

    @Post()
    book(@Body() offerID: string): string {
        Logger.log(`Booking request`);
        console.dir(offerID);

        // Get the Wishes associated to this offerID
        const initialWish = this.bookingsService.getWishesFromOfferID(offerID)

        this.bookingsService.startBookValidation(initialWish);
        return `Book validation initiated`;
    }

    @Post('payment')
    payBooking(@Body() wishPayment: WishPaymentDTO): string {
        // let paid = false //this.bookingsService.payAndBookById(bookingID)
        Logger.log("Wish received : " + wishPayment);
        return "OK";
    }
}
