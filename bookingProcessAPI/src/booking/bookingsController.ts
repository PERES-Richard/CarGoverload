import { Body, Controller, Get, Post, Logger } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { WishDTO } from 'src/models/wish_dto';
import { BookingsService } from "./bookings.service";

@Controller('booking-process/booking/')
export class BookingsController {
    constructor(private readonly bookingsService: BookingsService) {
    }
    @EventPattern("book-validation-result")
    bookValidationResultHandler(@Payload() data) {
        Logger.log(`Book validation result received. Processing validity...`);
        // TODO Check if offer selected is in the Offer's List
        // Unmarshall data
        let result
        this.bookingsService.handleBookValidation(result)
        // print result
    }

    @Post()
    book(@Body() offerID: string): string {
        Logger.log(`Booking request`);
        console.dir(offerID);

        // Get the Wishes associated to this offerID
        let initialWish = this.bookingsService.getWishesFromOfferID(offerID)

        this.bookingsService.startBookValidation(initialWish);
        return `Book validation initiated`;
    }

    @Get('/payment')
    payBooking(@Body() bookingID: string): string {
        let paid = this.bookingsService.payAndBookById(bookingID)

        if(!paid)
            return "Payment for booking n°" + bookingID + " failed.";

        return "Successfully Paid. Booking n°" + bookingID + " is complete.";
    }
}
