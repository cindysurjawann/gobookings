go build -o bookings.exe .
bookings.exe -dbname=gobookings -dbuser=postgres -dbpass=simaS123 -dbport=5431 -cache=false -production=false 