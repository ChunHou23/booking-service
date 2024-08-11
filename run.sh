#!/bin/bash

go build -o booking-service cmd/web/*.go
./booking-service -dbname=bookings -dbuser=chunhou -cache=false -production=false