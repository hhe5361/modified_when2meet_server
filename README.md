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


## Endpoints

### 1. Create Room
- **Method**: POST
- **URL**: `/rooms`
- **Description**: Creates a new meeting room
- **Request Body**:
```json
{
    "room_name": "string",
    "time_region": "string",
    "start_time": "number",
    "end_time": "number",
    "is_online": "boolean",
    "voteable_rooms": [
        {
            "year": "number",
            "month": "number",
            "day": "number"
        }
    ]
}
```
- **Response**:
```json
{
    "message": "Room created successfully",
    "data": {
        "url": "string"
    }
}
```

### 2. Register/Login
- **Method**: POST
- **URL**: `/rooms/:url/login`
- **Description**: Register a new user or login to an existing room
- **Request Body**:
```json
{
    "name": "string",
    "password": "string",
    "time_region": "string"
}
```
- **Response**:
```json
{
    "message": "Success",
    "data": {
        "user": {
            "id": "number",
            "name": "string",
            "time_region": "string",
            "available_time": []
        },
        "jwt_token": "string"
    }
}
```

### 3. Get Room Info
- **Method**: GET
- **URL**: `/rooms/:url`
- **Description**: Get room information and all users' details
- **Response**:
```json
{
    "message": "Success",
    "data": {
        "roomInfo": {
            "room": {
                "id": "number",
                "name": "string",
                "url": "string",
                "time_region": "string",
                "start_time": "number",
                "end_time": "number",
                "is_online": "boolean"
            },
            "dates": [
                {
                    "year": "number",
                    "month": "number",
                    "day": "number"
                }
            ]
        },
        "users": [
            {
                "user": {
                    "id": "number",
                    "name": "string",
                    "time_region": "string"
                },
                "available_time": [
                    {
                        "date": "string",
                        "hour_start_slot": "number",
                        "hour_end_slot": "number"
                    }
                ]
            }
        ]
    }
}
```

### 4. Get User Detail
- **Method**: GET
- **URL**: `/rooms/:url/user`
- **Auth Required**: Yes
- **Description**: Get current user's details and available times
- **Response**:
```json
{
    "message": "Success",
    "data": {
        "user": {
            "id": "number",
            "name": "string",
            "time_region": "string",
            "available_time": [
                {
                    "date": "string",
                    "hour_start_slot": "number",
                    "hour_end_slot": "number"
                }
            ]
        }
    }
}
```

### 5. Vote Time
- **Method**: PUT
- **URL**: `/rooms/:url/times`
- **Auth Required**: Yes
- **Description**: Update user's available time slots
- **Request Body**:
```json
{
    "date": "string (YYYY-MM-DD)",
    "hour_start_slot": "number",
    "hour_end_slot": "number"
}
```
- **Response**:
```json
{
    "message": "Vote time updated successfully"
}
```

### 6. Get Result
- **Method**: GET
- **URL**: `/rooms/:url/result`
- **Description**: Get meeting room results
- **Response**:
```json
{
    "message": "Success",
    "data": {
        "room": {
            "id": "number",
            "name": "string",
            "url": "string",
            "time_region": "string",
            "start_time": "number",
            "end_time": "number",
            "is_online": "boolean"
        },
        "dates": [
            {
                "year": "number",
                "month": "number",
                "day": "number"
            }
        ]
    }
}
```

## Error Responses
All endpoints may return the following error responses:
```json
{
    "error": "string"
}
```

Common status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error