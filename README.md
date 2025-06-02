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

# When2Meet API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Create Room
Creates a new meeting room with specified dates and time slots.

```http
POST /rooms
```

#### Request Body
```json
{
    "roomName": "Team Meeting",
    "timeRegion": "Asia/Seoul",
    "startTime": 9,
    "endTime": 18,
    "isOnline": false,
    "voteableRooms": [
        {
            "year": 2024,
            "month": 3,
            "day": 1
        }
    ]
}
```

#### Response
```json
{
    "message": "Room created successfully",
    "data": {
        "url": "room-abc123"
    }
}
```

### Register/Login User
Register a new user or login to an existing room.

```http
POST /rooms/:url/login
```

#### Request Body
```json
{
    "name": "John Doe",
    "password": "yourpassword",
    "timeRegion": "Asia/Seoul"
}
```

#### Response
```json
{
    "message": "Success",
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "timeRegion": "Asia/Seoul",
            "availableTimes": []
        },
        "jwt_token": "eyJhbGciOiJIUzI1NiIs..."
    }
}
```

### Get Room Information
Retrieves room details including all users and their available times.

```http
GET /rooms/:url
```

#### Response
```json
{
    "message": "Success",
    "data": {
        "roomInfo": {
            "id": 1,
            "name": "Team Meeting",
            "url": "room-abc123",
            "startTime": 9,
            "endTime": 18,
            "timeRegion": "Asia/Seoul",
            "isOnline": false
        },
        "vote_table": {
            "2024-03-01": [
                {
                    "hour": 9,
                    "users": ["John Doe", "Jane Smith"]
                }
            ]
        }
    }
}
```

### Get User Details
Retrieves details of the authenticated user.

```http
GET /rooms/:url/user
```

#### Headers
```
Authorization: Bearer <jwt_token>
```

#### Response
```json
{
    "message": "Success",
    "data": {
        "user": {
            "id": 1,
            "name": "John Doe",
            "timeRegion": "Asia/Seoul",
            "availableTimes": [
                {
                    "date": "2024-03-01",
                    "hourStartSlot": 9,
                    "hourEndSlot": 12
                }
            ]
        }
    }
}
```

### Vote Available Time
Record or update user's available time slots.

```http
PUT /rooms/:url/times
```

#### Headers
```
Authorization: Bearer <jwt_token>
```

#### Request Body
```json
{
    "date": "2024-03-01T00:00:00Z",
    "hourStartSlot": 9,
    "hourEndSlot": 12
}
```

#### Response
```json
{
    "message": "Vote time updated successfully"
}
```

## Error Responses
All endpoints may return the following error responses:

```json
{
    "error": "Error message description"
}
```

Common HTTP Status Codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

## Time Region Support
The API supports the following time regions:
- Asia/Seoul
- UTC
- America/New_York
- Europe/London 