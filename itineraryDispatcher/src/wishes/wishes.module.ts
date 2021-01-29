import {Module} from '@nestjs/common';
import {WishesController} from './wishes.controller';
import {WishesService} from './wishes.service';
import {ClientsModule, Transport} from "@nestjs/microservices";
import {KAFKA_HOST} from "../env_variable";

@Module({
  imports:[    ClientsModule.register([
    {
      name: 'DISPATCHER_SERVICE',
      transport: Transport.KAFKA,
      options: {
        client: {
          clientId: 'dispatcher',
          brokers: [KAFKA_HOST],
        },
        consumer: {
          groupId: 'dispatcher-consumer'
        }
      }
    },
  ]),],
  controllers: [WishesController],
  providers: [WishesService]
})
export class WishesModule {}
