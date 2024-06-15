# Chat System Microservice
*Assupmtions:*<br>
- Usernames are unique.
- Data fetching pattern is described similar to popular apps, such as Whatsapp, FB Messenger, Telegram, etc...
Where users once login can see a list of all messages sent to them or by them sorted by time in which they are sent in a DESC manner.

- Messages history retrieval does not apply grouping by sender or receiver yet!
- For logical & physical data modeling purposes, de-normalization of data schema is necessary to serve optimal performance. This leads to a bit of sacrifice when the need to write data could require batching sometimes.
Minimizing Batching/Atomic operations whenever possible to avoid impacting performance.

- Layer-first code organization approach as opposed to entity-first or feature-first for 2 reasons:<br>
  1. Avoiding cyclic imports.
  2. Making it familiar & straightforward.
However, yet some layers are ignored, again, for the sake of the demo and focusing on core functionality. One of such layers can be Repo to encapsulate data access logic,

*Disclaimers:*
- Since this is a development assignment meant for exploration & situation assessment. No database users/roles are created for the sake of the app usage & connection. Default credentials are used.
Set your own for security measures and best practices.

- No data seeding!

- Other environment parameters can be set & used across the app for more robust experience. Such as `DEBUG` & `APP_ENV` flags e.g

- Logging is subject for enhancement.
- Lack of Testing due to tight deadline in a holiday season.
- Simple input validation is conducted for the purpose of the demo. Rigorous validation with nicer error handling can be something to consider.

## How to Build and Run

### Prerequisites
- Docker
- Docker Compose

### Steps
1. Clone the repository.
2. Run `docker compose up --build`.
3. Ensure all services are up & healthy first.
4. Migrate the DataStore
Run the following command in CLI pointing to the project root. So, you can get your database schema created:
    ```bash
    $ migrate -path internal/db/migrations -database cassandra://cassandra:cassandra@localhost:9042/chat up
    ```
    <br>In case you encounter a db dirty state, just connect to Cassandra via any client to drop/remove all tables including the schema_migrations table under the chat keyspace.

    <br>
    SQLTools VSCode extension by `Matheus Teixeira` is a good one for GUI experience.

5. Access the service at `http://localhost/api/v1/<path>`.
You can utilize the included POSTMAN file with docs and environment setup to import at your end for convenience. "Find them at the root directory of the project in json format"

    Remember to register & login to grab the token and set its value in the Auth headers of upcoming requests to protected endpoints.
6. For Redis GUI client you can use `Redis Insight` as a good tool, works well on Linux/Ubuntu. Find your as per your respective platform though!<br> Or get savvy & jump right into the CLI mystical world!!

## API Endpoints
- `POST /register` - Register a new user
- `POST /login` - Login a user
- `POST /send` - Send a message
- `GET /messages` - Retrieve message history

## License
This is a free software distributed under the terms of the WTFPL license along with MIT license as dual-licensed, You can choose whatever works for you.<br/><br/>
Review the attached License file within the source code for mor details.

## Wanna Contribute?
Shout out!
