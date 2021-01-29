const { Kafka } = require('kafkajs')
const service = require('../services/carLocation.service')

let kafka
let ready = false

async function initConnection() {
    const HOST = process.env.KAFKA_HOST ? process.env.KAFKA_HOST : "kafka-service"
    const PORT = process.env.KAFKA_PORT ? process.env.KAFKA_PORT : "9092"

    kafka = new Kafka({
        clientId: 'car-location',
        brokers: [HOST+":"+PORT]
    })

    await subscribeToTopics().then(() => {
        ready = true;
    })
}

async function subscribeToTopics() {
    await createConsumer("car-location-new-search", "new-search", (message => service.newSearch(message.value)))
    await createConsumer("car-location-validation-search", "validation-search", (message => console.log(message)))
}

async function createConsumer(groupId, topic, callback) {
    const consumer = kafka.consumer({ groupId })

    await consumer.connect()
    await consumer.subscribe({ topic })

    await consumer.run({
        eachMessage: async ({ topic, partition, message }) => {
            callback(message)
        },
    })
}

async function sendMessage(topic, value) {
    let timeCount = 0
    const WAITING_TIME = 50000
    while (!ready && timeCount < WAITING_TIME) {
        timeCount++
    }
    if (timeCount >= WAITING_TIME) {
        throw 'Error: kafka bus not ready'
    }
    const producer = kafka.producer()

    await producer.connect()
    await producer.send({
        topic,
        messages: [
            { value },
        ],
    })

    await producer.disconnect()
}

module.exports = {
    initConnection,
    sendMessage
}
