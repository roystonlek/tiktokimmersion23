# TikTok Immersion Program 2023 Assignment

![Tests](https://github.com/TikTokTechImmersion/assignment_demo_2023/actions/workflows/test.yml/badge.svg)

## Royston Lek Attempt for TikTok Immersion Program 2023 Assignment:

### Requirements: 
- Docker ( https://docs.docker.com/get-docker/ ) 

### To run :
1) Start up docker 
2) In your terminal run `docker compose run --build -d` 
3) The API are ready at `127.0.0.1:8080` at 
- GET `/ping` 
- POST `/api/send` 
  - body : {
    "chat":chat,
    "text":text,
    "sender":sender
    }
- GET `/api/pull` 
  - body : {
    "chat":chat,
    "cursor":cursor,
    "limit":limit,
    "reverse":reverse
  }
