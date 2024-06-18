# Chat System Microservice
*Assumptions:*
- Usernames are unique.
- Data fetching pattern is described similar to popular apps, such as Whatsapp, FB Messenger, Telegram, etc...
Where users once login can see a list of all messages sent to them or by them sorted by time in which they are sent in a DESC manner.

- Messages history retrieval does not apply grouping by sender or receiver yet!

- For logical & physical data modeling purposes, de-normalization of data schema is necessary to serve optimal performance & promote efficient retrieval. This leads to a bit of sacrifice when the need to write data could require batching sometimes.<br>
Minimizing Batching/Atomic operations whenever possible to avoid impacting performance should be kept in mind though.

- Layer-first code organization approach as opposed to entity-first or feature-first for 2 reasons:<br>
  1. Avoiding cyclic imports.
  2. Making it familiar & straightforward.
However, yet some layers are ignored, again, for the sake of the demo and focusing on core functionality. One of such layers can be Repo to encapsulate data access logic,

- Token-based authentication is the method employed for Auth. Since a client-server arch can provide auth in different methods/mechanisms, this is the best fit here to avoid implications of other ways. Such other ways, can be listed as follows:
  - <b>Password-Based:</b> Where users are required to send credentials each time they access a protected resource. That's not efficient in our case.
  - <b>MFA:</b> Not suitable for the sake of this demo. As it adds extra layers that could lead to requiring more than just a password.
  - <b>OAuth:</b> We need to stay basic, not fancy.
  - <b>Stateful Authentication:</b> Why keep our server busy managing sessions and analyzing cookies while we can stay stateless, can't we?!<br>

  Thus JWT is picked up for easy interaction & smooth communication between the two parties (client & server) with pre-embedded credentials came from first auth operation as a handshake. Token expiration after one day would make sense to mitigate hacks, yet provide usage convenience.<br>
  <b>Note:</b> Important to understand that adopting HTTPS is the most important way to protect sniffing tokens over the network as it encrypts the data transmitted between the client & server or among services. Get rid of men in the middle in production settings.

*Disclaimers:*
- Since this is a development assignment meant for exploration & situation assessment. No database users/roles are created for the sake of the app usage & connection. Default credentials are used.
Set your own for security measures and best practices.

- No data seeding!

- Other environment parameters can be set & used across the app for more robust experience. Such as `DEBUG` & `APP_ENV` flags e.g

- Logging is subject for enhancement.
- Lack of Integration-Testing/E2E tests to verify API functionality. Due to tight deadline in a holiday season.
- Simple input validation is conducted for the purpose of the demo. Rigorous validation with nicer error handling can be something to consider.
- No Rate Limiting.

## How to Build and Run

### Prerequisites
- Docker
- Docker Compose

### Steps
1. Clone the repository.
2. Set the Environment Variables. Find the file named `.env.example` in the root directory of the project. Duplicate the file in the same path/location, rename the new one `.env` then set the values of the environment variables listed within the file according to your environment respectively.

    <br>Well, you may wanna keep `CASSANDRA_NODES` && `AUTH_HEADER_PREFIX` vars as-is with current values if you do wanna use the provided `docker-compose.yaml` file and `POSTMAN` collection & env without creating your own!

    Now you should be good to go.

    Unless you'd like to play a bit with `.air.toml` configurations. Consult the [Air's documentation]([URL](https://github.com/air-verse/air)) for more details/instructions in this regard.

3. Build:
    ```bash
    $ docker compose up --build
    ```
4. Ensure all services are up & healthy first.
5. Migrate the DataStore:<br>
    Install [golang-migrate/migrate CLI package]([URL](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)) to use.<br>
    Better with instructions for the respective OS you have as opposed to the "With Go toolchain" instructions on that page.<br>
    Additionally, [Docker usage]([URL](https://github.com/golang-migrate/migrate/?tab=readme-ov-file#docker-usage)) can be a breeze!

    Run the following command in CLI pointing to the project root. So, you can get your database schema created:<br>
    ```bash
    migrate -path internal/cassandra/migrations -database cassandra://cassandra:cassandra@localhost:9042/chat up
    ```
    <br>In case you encounter a db dirty state, just connect to Cassandra via any client to drop/remove all tables including the `schema_migrations` table under the `chat` keyspace. Or be a good citizen and force to the correct version instead. Then re-apply the migration direction desired!

    <br>I believe you know how to clear all db entries and drop the schema all at once:
    ```bash
    migrate -path internal/cassandra/migrations -database cassandra://cassandra:cassandra@localhost:9042/chat down
    ```
    Visiting [migrate's docs]([URL](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#usage)) is both helpful & recommended as usual.<br>

    <b>SQLTools</b> vsCode extension by <b>Matheus Teixeira</b> is a good one for GUI experience.

6. Access the service at `http://localhost/api/v1/<path>`.<br>
    You can utilize the included `POSTMAN` file with docs and environment setup to import at your end for convenience. "Find them at the root directory of the project in `json` format"

    Remember to register & login to grab the `token` and set its value in the Auth headers of subsequent requests to protected endpoints.<br>

7. For Redis GUI client you can use `Redis Insight` as a good tool, works well on Linux/Ubuntu. Find your as per your respective platform though!<br> Or get savvy & jump right into the CLI mystical world!!

## How to Test
First ensure Cassandra DB server is up & running on your localhost.
Or somewhere else but in this case you will need to adjust the tests cluster addresses manually in the code. Refer to `internal/services/test_utils.go` for some constants like `CLUSTER_ADDRS`.

#### Via IDE Comfy tools For GUI Folks
You can use the IDE comfy testing tools. However, please make sure to adjust the IDE settings for `Go: Test Timeout` to be set for something like `300s` which is equivalent to 5 minutes in such case. So, you allow it to take its time to establish a connection to cassandra DB and do all required jobs for test suites such as creating tables, seeding data, etc.
<br><br>
If you're on VSCode, adjusting timeout flag can be found as follows:
- File menu > Preferences > Settings. (Or hit CTRL + Comma keys on keyboard if you're on Linux)
- Search for `timeout`.<br>
If you've Go extension `Go Team at Google` `~v0.41.4` installed, then you should be able to locate that setting.
- Find that one setting under Extensions > Go.
- Set the value to nearly 300 seconds or more. As Cassandra connection seems very expensive operation.

#### Via CLI For Tech-Savvy
Visit the project-root directory from CLI and run:
```bash
$ go test ./... -v
```
For displaying of every single tiny log to stdout and stay on top of your game and be in control to check what went well and what goes wrong.
<br><br>
Run:
```sh
$ go test ./...
```
For no verbosity. If you're only interested in finding what passed and what failed.
<br>
Set/Pass some flags like `-timeout`to control tests timeout "defaults to 10 minutes". And `-run` to only run a particular test.

Head to `go help test` && `go help testflag` utilities docs for more detailed tips and hints on testing go modules. Good luck then!

## API Endpoints
- `POST /register` - Register a new user
- `POST /login` - Login a user
- `POST /send` - Send a message
- `GET /messages` - Retrieve message history

## License
This is a free software distributed under the terms of the `WTFPL` license along with MIT license as dual-licensed, You can choose whatever works for you.<br/><br/>
Review the attached `License` file within the source code for mor details.

## TODOS
Here's a list of bunch of stuff to be done when time allows!
- Integrate Prometheus & Grafana for metrics & monitoring.
- Write integrations test cases to ensure tests cover user registration, login, message sending, and message retrieval.

## Wanna Contribute?
Shout out!
