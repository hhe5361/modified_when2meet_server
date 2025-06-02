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
    "room_name": "Team Meeting",
    "time_region": "Asia/Seoul",
    "start_time": 9,
    "end_time": 18,
    "is_online": false,
    "voteable_rooms": [
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
    "time_region": "Asia/Seoul"
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
            "time_region": "Asia/Seoul",
            "available_times": []
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
            "start_time": 9,
            "end_time": 18,
            "time_region": "Asia/Seoul",
            "is_online": false,
            "created_at": "2024-03-01T00:00:00Z",
            "updated_at": "2024-03-01T00:00:00Z"
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
            "time_region": "Asia/Seoul",
            "available_times": [
                {
                    "id": 1,
                    "user_id": 1,
                    "date": "2024-03-01T00:00:00Z",
                    "hour_start_slot": 9,
                    "hour_end_slot": 12,
                    "created_at": "2024-03-01T00:00:00Z",
                    "updated_at": "2024-03-01T00:00:00Z"
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
    "hour_start_slot": 9,
    "hour_end_slot": 12
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