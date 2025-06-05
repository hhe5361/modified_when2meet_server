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
**Endpoint:** `POST /room`  
**Description:** Creates a new meeting room  
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
**Response:**
```json
{
  "message": "Room created successfully",
  "data": {
    "url": "string"
  },
  "error": null
}
```

### 2. Get Room Info
**Endpoint:** `GET /room/:url`  
**Description:** Retrieves room information including users and their available times  
**Response:**
```json
{
  "message": "Success",
  "data": {
    "roomInfo": {
      "id": "integer",
      "room_name": "string",
      "time_region": "string",
      "start_time": "integer",
      "end_time": "integer",
      "is_online": "boolean",
      "created_at": "datetime",
      "updated_at": "datetime"
    },
    "vote_table": {
      "users": [
        {
          "id": "integer",
          "name": "string",
          "time_region": "string",
          "available_times": [
            {
              "date": "datetime",
              "hour_start_slot": "integer",
              "hour_end_slot": "integer"
            }
          ]
        }
      ],
      "dates": [
        {
          "year": "integer",
          "month": "integer",
          "day": "integer"
        }
      ]
    }
  },
  "error": null
}
```

### 3. Register/Login
**Endpoint:** `POST /room/:url/register`  
**Description:** Registers a new user or logs in an existing user  
**Request Body:**
```json
{
  "name": "string",
  "password": "string"
}
```
**Response:**
```json
{
  "message": "Success",
  "data": {
    "user": {
      "id": "integer",
      "name": "string",
      "time_region": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    },
    "jwt_token": "string"
  },
  "error": null
}
```

### 4. Vote Time
**Endpoint:** `POST /room/:url/vote`  
**Description:** Updates user's available time slots (requires JWT authentication)  
**Request Body:**
```json
{
  "times": [
    {
      "date": "datetime",
      "hour_start_slot": "integer",
      "hour_end_slot": "integer"
    }
  ]
}
```
**Response:**
```json
{
  "message": "Vote time updated successfully",
  "data": null,
  "error": null
}
```

### 5. Get User Detail
**Endpoint:** `GET /user`  
**Description:** Retrieves authenticated user's details (requires JWT authentication)  
**Response:**
```json
{
  "message": "Success",
  "data": {
    "user": {
      "id": "integer",
      "name": "string",
      "time_region": "string",
      "created_at": "datetime",
      "updated_at": "datetime"
    }
  },
  "error": null
}
```

### 6. Get Result
**Endpoint:** `GET /room/:url/result`  
**Description:** Retrieves meeting result information  
**Response:**
```json
{
  "message": "Success",
  "data": {
    "room": {
      "id": "integer",
      "room_name": "string",
      "time_region": "string",
      "start_time": "integer",
      "end_time": "integer",
      "is_online": "boolean",
      "created_at": "datetime",
      "updated_at": "datetime"
    },
    "dates": [
      {
        "year": "integer",
        "month": "integer",
        "day": "integer"
      }
    ]
  },
  "error": null
}
```

### Error Responses
All endpoints may return the following error responses:
```json
{
  "message": "",
  "data": null,
  "error": "Error message"
}
```

### Authentication
- JWT token is required for protected endpoints
- Token should be included in the Authorization header
- Token contains user ID and room ID information

### Notes
- All datetime fields are in ISO 8601 format
- Time slots are represented as integers (e.g., 9 for 9:00 AM)
- The API uses snake_case for all field names
- All endpoints return a consistent response structure with message, data, and error fields

## Time Region Support
The API supports the following time regions:
- Asia/Seoul
- UTC
- America/New_York
- Europe/London 