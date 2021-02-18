import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { KAFKA, PORT } from "./env_variable";
import { Transport } from "@nestjs/microservices";

async function bootstrap() {
  const app = await NestFactory.create(AppModule, { cors: true });
  app.connectMicroservice({
    transport: Transport.KAFKA,
    options: {
      client: {
        brokers: [KAFKA],
      },
      consumer: {
        groupId: 'booking-process-consumer',
      }
    }
  });
  await app.startAllMicroservicesAsync();
  await app.listen(PORT);
}
bootstrap();
