import { Module } from '@nestjs/common';
import { BookingsController } from './bookingsController';
import { BookingsService } from './bookings.service';
import { ClientsModule, Transport } from "@nestjs/microservices";
import { KAFKA } from "../env_variable";

@Module({
  imports: [ClientsModule.register([
    {
      name: 'BOOKINGPROCESS_SERVICE',
      transport: Transport.KAFKA,
      options: {
        client: {
          clientId: 'bookingProcess',
          brokers: [KAFKA],
        },
        consumer: {
          groupId: 'booking-process-consumer'
        }
      }
    },
  ]),],
  controllers: [BookingsController],
  providers: [BookingsService]
})
export class BookingModule { }
