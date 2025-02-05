# Booking REST API Service
## Software Requirements
> [!NOTE]
> Versions specified are not strict requirements, but no guarantees can be made for unspecified versions.

*Required tools:*
- Go 1.23.*
- MySQL 8.4.3

## Local Development
Instructions for setting up the local development environment.

### TCP

1. Open your workspace in terminal and clone the project:
    - `git clone https://github.com/mnty4/booking`

2. Follow instructions to install MySQL 8.4.3 for your distribution:
https://dev.mysql.com/doc/refman/8.4/en/getting-mysql.html

> [!WARNING]
> Make sure to replace `<password>` with your own custom password.

3. Run the following commands in MySQL to setup the database and dev user:
    - `CREATE DATABASE IF NOT EXISTS booking;`
    - `CREATE USER IF NOT EXISTS dev_user@localhost IDENTIFY BY <password>;`
    - `GRANT ALL PRIVILEGES ON booking.* TO dev_user@localhost;`
    - `FLUSH PRIVILEGES;`

4. Run the following command in the terminal to populate the database from a sample backup:
    - `mysql -u dev_user -p booking < /PATH/TO/SRC/backup/booking.sql`

5. Follow instructions to install Go 1.23.*:
https://go.dev/doc/install

6. Copy example tcp env to `.env.tcp`:
    - `cat .env.tcp-example > .env.tcp`

7. Set the env variable, MYSQL_PASSWORD in `.env.tcp` to the same password you used when configuring your mysql dev user.

7. Run the following commands in the terminal from the project directory to build the binary and run the server:
    - `go build -o ./bin/bookingSvc ./cmd/bookingSvc`
    - `./bin/bookingSvc -env env.tcp`

8. Enter `Ctrl+C` to shutdown the server (the server should shutdown gracefully on signals SIGINT or SIGTERM).





