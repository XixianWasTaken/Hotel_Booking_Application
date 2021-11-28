go build -o bookings cmd/web/*.go

#host = localhost port = 5432 dbname = bookings user = postgres password = 123456"
./bookings -dbhost=localhost dbport=5432 -dbname=booking -dbuser=postgres -dbpass=123456 -cache=false -production=false