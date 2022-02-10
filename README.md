# Cloud Native todolist backend

Today we have a lot of possibilities to bootstrap a application or build minimal product using Cloud Services like [Supabase](https://supabase.com) and [Ory Cloud](https://www.ory.sh/cloud/). This demo app shows an example of simple todo list backend API.

1. Authentication/User login/Forgot password are handled by [Ory Cloud](https://www.ory.sh/cloud/) or [Ory Kratos](https://www.ory.sh/kratos)
2. Database layer is implemented by [Supabase](https://supabase.com) database

## Authentication

For the authentication part you can do two things:

1. Use Ory Cloud (just follow the onboarding process to their Cloud services)
2. Use self-hosted Ory Kratos installation (available through docker-compose in this project)

## Authentication Flow

TODO:
- [ ] Add Mermaid diagram
- [ ] Add description 


## Database
Currently it's a showcase for Supabase, but the roadmap has the following parts

- [x] Supabase Database
- [ ] Firebase Database
- [ ] Appwriter Database
- [ ] Postgres/MySQL database with examples using RDS 
- [ ] Create an issue to add a new database example

## API

API is implemented using the following toolchain

1. buf.build to improve development experience working with protobufs
2. gRPC as a main technology to build API
3. gRPC gateway to add support of RESTful API for browser

## How does authentication works for REST APIs

You can check out the code of KratosMiddleware located in `middleware` package. The process is simple and it's used to the browsers only (mobile and desktop)

[![](https://mermaid.ink/img/pako:eNptUsFKJDEQ_ZUiJwUXB3Q85CDIrMiwF1HmtL1okdT0hOlOepNqZFaE_QEPfoG3RW9-k1_gJ1iZtCs4U5eqJO-9elXkVplgSWmV6HdP3tB3h3XEtvIg0WFkZ1yHnqFPFAETzCQn2Hl7fHh6_Xv_9vjvRcrnodzd5C0jckizaeb-WNcghx19eDgeb4EPkP_gDDw42AKsY2cyrL44n8AZMt3gSmTHo9FomyxxWIvmnDWPjgRVcHmyb8fHdfxpfmmYYNNIfw8n59Pyjg3Dpas9TP1-zrOu3OfILtbczmg4I4brEFdXZearRCm54K_BhLB0tMEqExaes-TZ8Qo4ouME8xha2L9ZBGwdkLddcJ6LAjWJAHteZIaRye1w7-2n4RNjpDtMgucYms3WsgcNZkFmCVigshiRhEVoqcOavvQK0f3ZaFTE8v40XBD30SeZ1bMYSxDmeYVfvH8ayaH2VEuxRWflB97mt0qJh5YqpaW0GJeVqvyd4PrOyqSn1nGISs9RfO0pMRYuV94ozbGnD9DwhQfU3TsYJvnK)](https://mermaid.live/edit#pako:eNptUsFKJDEQ_ZUiJwUXB3Q85CDIrMiwF1HmtL1okdT0hOlOepNqZFaE_QEPfoG3RW9-k1_gJ1iZtCs4U5eqJO-9elXkVplgSWmV6HdP3tB3h3XEtvIg0WFkZ1yHnqFPFAETzCQn2Hl7fHh6_Xv_9vjvRcrnodzd5C0jckizaeb-WNcghx19eDgeb4EPkP_gDDw42AKsY2cyrL44n8AZMt3gSmTHo9FomyxxWIvmnDWPjgRVcHmyb8fHdfxpfmmYYNNIfw8n59Pyjg3Dpas9TP1-zrOu3OfILtbczmg4I4brEFdXZearRCm54K_BhLB0tMEqExaes-TZ8Qo4ouME8xha2L9ZBGwdkLddcJ6LAjWJAHteZIaRye1w7-2n4RNjpDtMgucYms3WsgcNZkFmCVigshiRhEVoqcOavvQK0f3ZaFTE8v40XBD30SeZ1bMYSxDmeYVfvH8ayaH2VEuxRWflB97mt0qJh5YqpaW0GJeVqvyd4PrOyqSn1nGISs9RfO0pMRYuV94ozbGnD9DwhQfU3TsYJvnK)

## How does authentication works for gRPC APIs

Not implemented yet

## Installation and demo

```
git clone git@github.com:gen1us2k/cloudnative_todo_list
cd cloudnative_todo_list
docker-compose up -d
go run cmd/todolist/main.go
```

Open http://127.0.0.1:8081/api/todo to start demo

## Development

The project uses Go and Make for local development
```
build_grpc                     Generate files for gRPC
deps                           install binaries
lint                           Run lint against Go code
test                           Run tests
```

## Configuration
Configuration can be changed with environment variables
