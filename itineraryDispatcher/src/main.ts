import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import {KAFKA_HOST, PORT} from "./env_variable";
import {Transport} from "@nestjs/microservices";

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  app.connectMicroservice({
    transport: Transport.KAFKA,
    options: {
      client: {
        brokers: [KAFKA_HOST],
      },
      consumer: {
        groupId: 'dispatcher-consumer',
      }
    }});
  await app.startAllMicroservicesAsync();
  await app.listen(PORT);
}
bootstrap();
