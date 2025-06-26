# CodeBase -- interactive learning platform for building stuff 
--- 

## Motivation 

I want to build an interactive learning platform as it supposed to be -- fun and intuitive. This platform should be rather seen as a tool for learning how to code, rather than a 100% guarantee that you land a job after completing the course (which is bullshit, don't give your money to those) 

The problems I faced while learning and hopping from one source / tutorial / course to another were: 

- There are tons of info, so the feeling of constant overwhelming and being lost is familiar to me 
- I don't like videos for learning something practical -- it's better to learn by doing 
- EdTech platforms resemble each other -- your complete a task, you gain points, yayaya. Maybe I could think of more fun and sticky way to learn -- learning in a full survival mode 

## Stack 

- Microservices -- because it seems more intuitive to me 


## Root Structure

```
/codebase
├── services
│   ├── api-gateway
│   ├── user-service
│   ├── task-service
│   ├── missions-service
│   ├── xp-service
│   └── notification-service
├── pkg
│   └── common-lib
│       ├── config
│       ├── logger
│       └── error
├── infra
│   ├── docker-compose.yml
│   ├── k8s
│   └── helm
├── scripts
│   ├── migrate.sh
│   └── deploy.sh
├── ci
│   ├── github-actions.yml
│   └── lint.yml
└── README.md
```

---

## Service Layout (example: task-service)

```
services/task-service/
├── cmd
│   └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── handler
│   │   └── task_handler.go
│   ├── service
│   │   └── task_service.go
│   ├── repository
│   │   ├── repository.go
│   │   └── postgres
│   │       └── task_postgres.go
│   ├── model
│   │   └── task.go
│   └── logger
│       └── logger.go
├── migrations
│   └── 001_create_tasks_table.up.sql
├── Dockerfile
├── go.mod
└── go.sum
```

## Task Matrix 

| Module                 | Status         | Notes                          |
| ---------------------- | -------------- | ------------------------------ |
| api-gateway            | ✅ Completed    | Basic routing & proxy          |
| auth-service           | ⚠️  In Progress  | Handles JWT, login, refresh    |
| user-service           | ⚠️ In Progress | Add RBAC, 2FA                  |
| task-service           | ✅ Completed    | CRUD + filtering               |
| missions-service       | ⚠️ In Progress | Refactor to new service layout |
| xp-service             | ❌ Not Started  | Kafka consumer, Redis          |
| notification-service   | ❌ Not Started  | Email + webhook support        |
| Shared `common-lib`    | ❌ Not Started  | logger, config loader          |
| infra (docker-compose) | ⚠️ In Progress | define all services            |
| CI/CD                  | ⚠️ In Progress | GitHub Actions, lint           |
