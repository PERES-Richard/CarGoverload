#!/usr/bin/env bash

echo Test booking process service..
./healthcheck_service.sh http://localhost/booking-process/ok

echo Test car availability service..
./healthcheck_service.sh http://localhost/car-availability/ok

echo Test car booking service..
./healthcheck_service.sh http://localhost/car-booking/ok

echo Test car searching service..
./healthcheck_service.sh http://localhost/car-searching/ok
