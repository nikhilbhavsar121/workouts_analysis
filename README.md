# workouts_analysis

<!-- Install golang migrator -->
<!-- 
Step to follows
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate 
-->
Create the database with name workoutdb

Run the migration
migrate -path db/migration -database  "mysql://Username:password@tcp(127.0.0.1:3306)/workoutdb" -verbose up

set .env file
set the test db configuration in the test cases