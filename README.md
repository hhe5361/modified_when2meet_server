I made Better When2Meet
this repo is when2meet server

current when2Meet features
- make room
- login to specific room
- get specific room data
- post specific room data

I append below features on my project
- where to meet (not just when to meet).
- Time synchronization can be difficult when team members live in different countries, making it hard to decide on a time for online meetings.
- result visualization


ref project: https://www.when2meet.com/

# When2Meet Server

## API Documentation

### 1. Create Room
**Endpoint:** `POST /rooms`  
**Request Body:**
```json
{
  "room_name": "string",
  "time_region": "string",
  "start_time": "integer",
  "end_time": "integer",
  "is_online": "boolean",
  "voteable_rooms": [
    {
      "year": "integer",
      "month": "integer",
      "day": "integer"
    }
  ]
}
```

### 2. Login
**Endpoint:** `POST /rooms/:url/login`  
**Request Body:**
```json
{
  "name": "string",
  "password": "string",
  "time_region": "string"
}
```

### 3. Vote Time
**Endpoint:** `PUT /rooms/:url/times`  
**Request Body:**
```json
{
  "times": [
    {
      "date": "YYYY-MM-DD",
      "hour_start_slot": "integer",
      "hour_end_slot": "integer"
    }
  ]
}
```

### 4. Get Room Info
**Endpoint:** `GET /rooms/:url`  
**Response:**
```json
{
  "message": "Success",
  "data": {
    "room": {
      "id": "integer",
      "name": "string",
      "time_region": "string",
      "start_time": "integer",
      "end_time": "integer",
      "is_online": "boolean"
    },
    "dates": [
      {
        "year": "integer",
        "month": "integer",
        "day": "integer"
      }
    ],
    "vote_table": {
      "YYYY-MM-DD": [
        {
          "hour": "integer",
          "users": ["string"]
        }
      ]
    }
  }
}
```
```

Additional suggestions:
1. Add error response examples for common cases
2. Document the JWT token format and required claims
3. Add rate limiting information if applicable
4. Add pagination details if applicable
5. Document the time slot format more clearly (e.g., 9 = 9:00 AM)
6. Add examples of successful and error responses for each endpoint

Would you like me to provide a complete corrected version of the API documentation?