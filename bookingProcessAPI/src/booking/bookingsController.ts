import { Body, Controller, Post, Logger } from '@nestjs/common';
import { EventPattern, Payload } from "@nestjs/microservices";
import { BookingsService } from "./bookings.service";
import {WishPaymentDTO} from "../models/wish_payment_dto";

@Controller('booking-process/booking/')
export class BookingsController {
    constructor(private readonly bookingsService: BookingsService) {
    }
    
    @EventPattern("book-confirmation")
    bookConfirmationResultHandler(@Payload() data) {
        Logger.log(`Book validation result received. Processing validity...`);
        this.bookingsService.handleBookConfirmation(data.value)
    }

    @Post('payment')
    payBooking(@Body() wishPayment: WishPaymentDTO): string {
        Logger.log("Wish received : ");
        console.dir(wishPayment)
        this.bookingsService.payAndBookById(wishPayment);
        return "Waiting for payment";
    }
}
