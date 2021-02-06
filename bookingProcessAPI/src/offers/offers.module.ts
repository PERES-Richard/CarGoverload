import { Module } from '@nestjs/common';
import { OffersController } from './offers.controller';
import { OffersService } from './offers.service';
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
  controllers: [OffersController],
  providers: [OffersService]
})
export class OffersModule { }
