import { Module } from '@nestjs/common';
import * as redis from 'redis';
import { REDIS_HOST } from '../env_variable';

const RedisProvider = {
    provide: "redis",
    useFactory: () => redis.createClient(`redis://${REDIS_HOST}`)
}

@Module({
    exports: ["redis"],
    providers: [RedisProvider]
})
export class RedisModule { }
