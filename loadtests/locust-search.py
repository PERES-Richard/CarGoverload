import random
import time
from locust import HttpUser, task, between


class QuickstartUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        user = {
            "phone": "0"+str(time.time())+str(random.randint(1, 180)),
        }

        search1 = {"carType": "Liquid", "numberOfCars": 2, "departureNode": "Marseille",
                   "arrivalNode": "Paris", "dateDeparture": "2006-01-02T15:04:05Z", }
        search = []
        search.append(search1)
        self.client.post("/booking-process/offers", json=search,
                         headers={"Content-Type": "application/json", "Origin": '*'})
