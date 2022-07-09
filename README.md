# FoodLogiQ demo GoLang REST Project

## Project Requirements
Pulled directly from [this doc](https://docs.google.com/document/d/1iHz8jhM0TINu6EquXMtnP7pvU2R2ErbLM8NilSIa7xk/edit)
and referenced in [stories file](./stories.txt).

## API Docs
All requests require an `Authorization: Bearer <token>` in order to be authenticated

| Path                               | Method | Input                                                                                                                                | Response Codes                                                                                                                  | Responses                                                                                                                                         |
|------------------------------------|--------|--------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------|
| https://host/event?event_id=string | GET    | A string ID for the event                                                                                                            | 200 - OK<br>400 - Request error <br>401 - Unauthorized<br>500 - MongoDB error                                                   | {"type":"shipping","contents":[{"gtin":"1234","lot":"adffda","bestByDate":"2021-01-13","expirationDate":"2021-01-17"}]}                           |
| https://host/event                 | POST   | {"type":"<string>","contents":[{"gtin":"<string>","lot":"<string>","bestByDate":"[Date in RFC3339]","expirationDate":"2021-01-17"}]} | 201-  Created<br>400 - Request malformed (with error details)<br>401 - Unauthorized<br>500 - MongoDB error                      | {"event_id": "<string>"}                                                                                                                          |
| https://host/event                 | DELETE | {"id":"<string">}                                                                                                                    | 200 - Success (deleted)<br>400 - Bad request (cannot match eventId and created_by)<br>401 - Unauthorized<br>500 - MongoDB Error | N/A                                                                                                                                               |
| https://host/event/all             | GET    | N/A                                                                                                                                  | 200 - Success<br>400 - No events found<br>401 - Unauthorized<br>500 - MongoDB Error                                             | Array of Events <br><br>[<br>{"type":"shipping","contents":[{"gtin":"1234","lot":"adffda","bestByDate":"2021-01-13","expirationDate":"2021-01-17] |

## How to build
1. Ensure `make` is installed on your system, along with `docker` and `docker-compose`
2. Run the command `make run-demo`
   1. This will run `make compose` & `make run-demo` which combined will:
      1. Build the app using docker compose
      2. Run the app in docker and expose ports
      3. Run the test `main.go` app which performs a create, get, list, and delete REST request
3. Additional commands can be seen by doing `make`

## Considerations
`Tokens` - There were a set of tokens as defined [here](./user.json). There are 2 options: load them into a DB, or load into memory and do mapping.
I chose to do in-memory mapping rather than placing in the DB for simplicity.

`Unit Testing` - Using the ginkgo and gomega framework is different from usual baked-in testing libraries used by Go. 
It allows for a more rich experience and pathing / asserting system than standard, and is a personal preference.
Converting tests to use standard go libraries can be done if required

`DB Testing` - The go library [mockery](https://github.com/vektra/mockery) is used to stub out the interface and handle mock responses.
There is no integration testing confirming the BSON requests.

`Returning the same error for GET /event/` - Per user story 3 (As an API consumer, I want a REST API to retrieve a specific event), there were 3 requirements for returning an error response:
1. If the ID provided is not found, return an appropriate failure status code
2. If the ID provided is found, but was not created by the user who is making the REST API call, return the same failure as if the ID was not found
3. If the event of the ID has been deleted, return the same failure as if the ID was not found

The nature of this AC seems to imply a SELECT operation, then a filter on the select. I have elected to not do that, and simply apply bson filtering. This is a performance improvement.

`Returning the same error for DELETE /event/` - Per user story 2 (As an API consumer, I want a REST API to delete a specific event)
1. If the supplied event ID does not exist, return an appropriate failure status code
2. If the supplied event ID does exist, but not accessible by the user, return an appropriate failure status code

Again, this seems to imply SELECT from Mongo, then filter by `created_by` field to then return an error stating the ID exists, but they don't own it.
I elected to return the same error and use the same bson filter. This not only is a performance increase, but a security improvement.
While threats are low considering what data is stored in an event, lets say an attacker KNOWS theres a specific event_id, they could more easily brute force a change if they only have to guess 1 field VS 2.


