# API Documentation

This document provides comprehensive documentation for all API endpoints in the Aetrix Observer backend.

## Base URL

```
http://<host>:<port>/api
```

---

## Table of Contents

1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Machine Management](#machine-management)
4. [Command Workflow Guide](#command-workflow-guide)
5. [Command & Container Operations](#command--container-operations)
6. [WebSocket Endpoints](#websocket-endpoints)

---

## Authentication

### 1. User Signup

**Endpoint:** `POST /api/auth/user/signup`

**Description:** Register a new user account.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)",
  "email": "string (required, valid email format)"
}
```

**Success Response (201):**
```json
{
  "message": "User created successfully",
  "user_id": 1
}
```

**Error Response (400):**
```json
{
  "error": "validation error message"
}
```

---

### 2. User Login

**Endpoint:** `POST /api/auth/user/login`

**Description:** Authenticate user and receive JWT token.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Success Response (200):**
```json
{
  "message": "Login successful",
  "user_id": 1,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- **400:** Invalid request format
  ```json
  {
    "error": "validation error message"
  }
  ```
- **401:** Invalid credentials
  ```json
  {
    "error": "Invalid username or password"
  }
  ```
- **500:** Token generation error
  ```json
  {
    "error": "Failed to generate token"
  }
  ```

---

## User Management

### 3. Get User Profile

**Endpoint:** `GET /api/user/profile`

**Description:** Retrieve authenticated user's profile information.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com"
}
```

**Error Responses:**
- **401:** Unauthorized
  ```json
  {
    "error": "Unauthorized"
  }
  ```
- **404:** User not found
  ```json
  {
    "error": "User not found"
  }
  ```

---

### 4. Update User Details

**Endpoint:** `PUT /api/user/update`

**Description:** Update user's username and/or email.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "username": "string (optional)",
  "email": "string (optional, valid email format)"
}
```

**Success Response (200):**
```json
{
  "message": "User updated successfully"
}
```

**Error Responses:**
- **400:** Invalid request format
- **401:** Unauthorized
- **404:** User not found

---

### 5. Change Password

**Endpoint:** `PUT /api/user/change-password`

**Description:** Change user's password.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "old_password": "string (required)",
  "new_password": "string (required)"
}
```

**Success Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Error Responses:**
- **400:** Invalid request or incorrect old password
  ```json
  {
    "error": "Old password is incorrect"
  }
  ```
- **401:** Unauthorized
- **404:** User not found

---

## Machine Management

### 6. Register Machine

**Endpoint:** `POST /api/user/machine`

**Description:** Register a new machine/agent.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "name": "string (required)",
  "users": [1, 2, 3], // Array of user IDs (required)
  "ip": "string (required)" // Machine IP address
}
```

**Success Response (200):**
```json
{
  "message": "machine created successfully",
  "id": 1
}
```

**Error Responses:**
- **400:** Invalid request
  ```json
  {
    "message": "invalid data",
    "error": "validation error details"
  }
  ```
- **500:** Database error

---

### 7. Get Machine Details

**Endpoint:** `GET /api/user/machine/:machine_id`

**Description:** Get detailed information about a specific machine.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `machine_id` (string): Machine ID

**Success Response (200):**
```json
{
  "message": "data retrieved",
  "data": {
    "id": 1,
    "name": "Production Server",
    "ip": "192.168.1.100",
    "creator_id": 1,
    "creator": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com"
    },
    "users": [
      {
        "id": 2,
        "username": "jane_doe",
        "email": "jane@example.com"
      }
    ]
  }
}
```

**Error Response (400):**
```json
{
  "message": "error while fetching data",
  "err": "error details"
}
```

---

### 8. List Owned Machines

**Endpoint:** `GET /api/user/machine/owned`

**Description:** List all machines created by the authenticated user.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "machines": [
    {
      "id": 1,
      "name": "Production Server",
      "ip": "192.168.1.100",
      "creator_id": 1,
      "users": [...]
    }
  ]
}
```

**Error Response (400):**
```json
{
  "message": "failed to fetch",
  "machines": []
}
```

---

### 9. List Usable Machines

**Endpoint:** `GET /api/user/machine/usable`

**Description:** List all machines the authenticated user has access to (including shared machines).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Success Response (200):**
```json
{
  "machines": [
    {
      "id": 1,
      "name": "Production Server",
      "ip": "192.168.1.100",
      "creator_id": 2,
      "users": [...]
    }
  ]
}
```

**Error Response (400):**
```json
{
  "message": "error while fetching data",
  "err": "error details"
}
```

---

### 10. Update Machine

**Endpoint:** `PUT /api/user/machine`

**Description:** Update machine information (name and/or IP). Only the machine creator can update.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "id": 1, // required
  "name": "string (optional)",
  "ip": "string (optional)"
}
```

**Success Response (200):**
```json
{
  "message": "machine updated",
  "machine": {
    "id": 1,
    "name": "Updated Name",
    "ip": "192.168.1.101"
  }
}
```

**Error Responses:**
- **400:** Invalid request, permission denied, or update failed
  ```json
  {
    "message": "No permission",
    "err": "Creator Mismatch"
  }
  ```

---

### 11. Add Users to Machine

**Endpoint:** `PUT /api/user/machine/add/user`

**Description:** Grant additional users access to a machine. Only the machine creator can add users.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "id": 1, // Machine ID (required)
  "users": [3, 4, 5] // Array of user IDs to add (required)
}
```

**Success Response (200):**
```json
{
  "message": "users added successfully",
  "data": {
    "id": 1,
    "name": "Production Server",
    "users": [...]
  }
}
```

**Error Responses:**
- **400:** Invalid request or unauthorized
- **500:** User not found or database error

---

### 12. Remove Users from Machine

**Endpoint:** `PUT /api/user/machine/remove/user`

**Description:** Revoke user access from a machine. Only the machine creator can remove users.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "id": 1, // Machine ID (required)
  "users": [3, 4] // Array of user IDs to remove (required)
}
```

**Success Response (200):**
```json
{
  "message": "users added successfully",
  "data": {
    "id": 1,
    "name": "Production Server",
    "users": [...]
  }
}
```

**Error Responses:**
- **400:** Invalid request or unauthorized
- **500:** Database error

---

### 13. Delete Machine

**Endpoint:** `DELETE /api/user/machine/:machine_id`

**Description:** Delete a machine. Only the machine creator can delete.

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `machine_id` (string): Machine ID

**Success Response (200):**
```json
{
  "message": "machine deleted successfully"
}
```

**Error Responses:**
- **400:** Invalid ID or deletion failed
  ```json
  {
    "message": "invalid credentials",
    "err": "error details"
  }
  ```

---

## Command Workflow Guide

This section describes the end-to-end flow for sending Docker commands: authentication, agent setup on a VM, and command payload formats.

### 1) Authenticate (Get JWT)

1. Sign up using `POST /api/auth/user/signup`
2. Log in using `POST /api/auth/user/login`
3. Copy the returned JWT token and send it in the `Authorization` header for all command requests:

```
Authorization: Bearer <jwt_token>
```

### 2) Register a Machine (Optional but Recommended)

Use `POST /api/user/machine` to create a machine record for UI/ownership tracking. The `machine_id` used in command requests must match the agent identifier used by the VM agent connection. Keep them consistent across your setup.

### 3) Run the Agent on a VM

On each VM that should execute Docker commands:

1. Install and start Docker Engine.
2. Ensure the agent process can access the Docker socket (e.g., user has permission to `/var/run/docker.sock`).
3. Configure the agent connection in `cmd/agent/main.go`:
   - `host`: the API server host and port (for example, `api.example.com:8080`)
   - `path`: keep as `/agent`
   - `clientName`: the agent ID (this must match `machine_id` in command requests)
4. Build and run the agent binary on the VM.

When the agent starts, it connects to:

```
ws://<host>/agent/<clientName>
```

### 4) Send Docker Commands (HTTP or WebSocket)

Use `POST /api/user/cmd` for request/response commands or `GET /api/user/cmd` (WebSocket) for streaming.

#### Command Types (Exact Values)

These are the accepted `command_type` values:

- `list_container`
- `start_new_container`
- `start_container`
- `stop_container`
- `restart_container`
- `delete_container`

#### Params Shape

Use the following `params` object shape:

```json
{
  "list": {
    "all": true,
    "size": true
  },
  "start": {
    "name": "my-nginx",
    "image": "nginx:latest",
    "host_port": 8080,
    "container_port": 80,
    "protocol": "tcp",
    "host_ip": "0.0.0.0"
  },
  "delete": {
    "force": true,
    "volume": true,
    "links": false
  }
}
```

Only include the relevant sub-object for the command you are sending.

#### Example: List Containers

```json
{
  "machine_id": "agent-1",
  "command_type": "list_container",
  "params": {
    "list": {
      "all": true,
      "size": true
    }
  }
}
```

#### Example: Start New Container

```json
{
  "machine_id": "agent-1",
  "command_type": "start_new_container",
  "params": {
    "start": {
      "name": "my-nginx",
      "image": "nginx:latest",
      "host_port": 8080,
      "container_port": 80,
      "protocol": "tcp",
      "host_ip": "0.0.0.0"
    }
  }
}
```

#### Example: Stop Container

```json
{
  "machine_id": "agent-1",
  "container_id": "abc123def456",
  "command_type": "stop_container"
}
```

### 5) Docker CLI Equivalents

The commands map to Docker CLI operations as follows:

- `list_container` → `docker ps` (use `-a` for `all=true`, `--size` for `size=true`)
- `start_new_container` → `docker run -d --name <name> -p <host_port>:<container_port>/<protocol> <image>`
- `start_container` → `docker start <container_id>`
- `stop_container` → `docker stop <container_id>`
- `restart_container` → `docker restart <container_id>`
- `delete_container` → `docker rm <container_id>` (add `-f` for `force=true`, `-v` for `volume=true`, `--link` for `links=true`)

---

## Command & Container Operations

### 14. Send Command (HTTP)

**Endpoint:** `POST /api/user/cmd`

**Description:** Send a Docker command to a machine and wait for the response (HTTP request-response pattern).

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Request Body:**
```json
{
  "machine_id": "string (required)",
  "container_id": "string (optional, required for container-specific commands)",
  "command_type": "string (required)",
  "params": {
    // Command-specific parameters
  }
}
```

**Command Types:**
- `list_container` - List all containers
- `start_new_container` - Create and start a new container
- `start_container` - Start an existing container
- `stop_container` - Stop a running container
- `restart_container` - Restart a container
- `delete_container` - Remove a container

**Example: List Containers**
```json
{
  "machine_id": "agent-1",
  "command_type": "list_container",
  "params": {
    "list": {
      "all": true,
      "size": true
    }
  }
}
```

**Example: Start New Container**
```json
{
  "machine_id": "agent-1",
  "command_type": "start_new_container",
  "params": {
    "start": {
      "name": "my-nginx",
      "image": "nginx:latest",
      "host_port": 8080,
      "container_port": 80,
      "protocol": "tcp",
      "host_ip": "0.0.0.0"
    }
  }
}
```

**Example: Stop Container**
```json
{
  "machine_id": "agent-1",
  "container_id": "abc123def456",
  "command_type": "stop_container"
}
```

**Example: Delete Container**
```json
{
  "machine_id": "agent-1",
  "container_id": "abc123def456",
  "command_type": "delete_container",
  "params": {
    "delete": {
      "force": true,
      "volume": true,
      "links": false
    }
  }
}
```

**Success Response (200):**
```json
{
  "data": "Container started successfully with ID: abc123def456",
  "error": null
}
```

**Error Responses:**
- **400:** Invalid request
  ```json
  {
    "error": "validation error message"
  }
  ```
- **500:** Command failed or timeout
  ```json
  {
    "error": "timeout waiting for response from machine agent-1"
  }
  ```

**Note:** This endpoint has a 30-second timeout for receiving responses from the agent.

---

### 15. Send Command (WebSocket)

**Endpoint:** `GET /api/user/cmd` (WebSocket upgrade)

**Description:** Send a command via WebSocket and receive streaming responses. Useful for long-running operations or when you need real-time updates.

**Headers:**
```
Authorization: Bearer <jwt_token>
Connection: Upgrade
Upgrade: websocket
```

**Query Parameters:**
- `machine_id` (string, required): Target machine ID
- `container_id` (string, optional): Container ID for container-specific commands
- `command_type` (string, required): Command type (same as HTTP endpoint)
- `params` (string, optional): URL-encoded JSON string of command parameters

**Example URL:**
```
ws://localhost:8080/api/user/cmd?machine_id=agent-1&command_type=list_container&params=%7B%22list%22%3A%7B%22all%22%3Atrue%7D%7D
```

**WebSocket Response Format:**
Messages are sent as JSON via the WebSocket connection:
```json
{
  "data": { /* response data */ },
  "error": null
}
```

**Connection Flow:**
1. Client upgrades HTTP connection to WebSocket
2. Server adds client as subscriber for the machine
3. Command is sent to the agent
4. Responses/streams are forwarded to the client via WebSocket
5. Connection closes after completion or timeout (30 seconds)
6. Client is automatically unsubscribed

**Error Response (400):**
```json
{
  "error": "invalid parameters"
}
```

**Error Response (500):**
```json
{
  "error": "Failed to upgrade to websocket"
}
```

---

## WebSocket Endpoints

### 16. Agent WebSocket Connection

**Endpoint:** `GET /agent/:machine_id` (WebSocket upgrade)

**Description:** WebSocket endpoint for agents to connect and receive commands. This is used by the agent service, not by end users.

**Path Parameters:**
- `machine_id` (string): Unique machine/agent identifier

**Connection Flow:**
1. Agent connects via WebSocket
2. Server registers the agent connection
3. Server can send commands to the agent
4. Agent sends responses and events back to server
5. Connection remains open for bidirectional communication

**Message Types from Agent:**

**1. Response Message:**
```json
{
  "type": "response",
  "machine_id": "agent-1",
  "payload": {
    "data": { /* command result */ },
    "error": null
  }
}
```

**2. Event Message:**
```json
{
  "type": "event",
  "machine_id": "agent-1",
  "payload": {
    "event": "UNEXPECTED_STOP",
    "container_id": "abc123",
    "data": { /* event details */ }
  }
}
```

**Event Types:**
- `UNEXPECTED_STOP` - Container stopped unexpectedly (triggers automatic restart with retry logic)

**3. Stream Message (Metrics):**
```json
{
  "type": "stream",
  "machine_id": "agent-1",
  "payload": {
    "data": {
      "cpu": 12.3,
      "memory": 512
    }
  }
}
```

**Message Types to Agent:**

**Command Message:**
```json
{
  "cmd": "start_container",
  "machine_id": "agent-1",
  "container_id": "abc123",
  "params": {
    // command-specific parameters
  }
}
```

---

## Error Handling

### Standard Error Response Format

All endpoints return errors in a consistent format:

```json
{
  "error": "Error message describing what went wrong"
}
```

Or:

```json
{
  "message": "Contextual message",
  "err": "Detailed error information"
}
```

### Common HTTP Status Codes

- **200 OK** - Successful operation
- **201 Created** - Resource created successfully
- **400 Bad Request** - Invalid request data or parameters
- **401 Unauthorized** - Missing or invalid authentication
- **404 Not Found** - Resource not found
- **500 Internal Server Error** - Server-side error

---

## Authentication & Authorization

### JWT Token

Most endpoints require authentication via JWT token in the `Authorization` header:

```
Authorization: Bearer <jwt_token>
```

The token is obtained from the login endpoint and contains the user ID.

### Middleware

- **User Middleware**: Validates JWT token and extracts `user_id`
- **Machine Middleware**: Validates machine-specific permissions (currently not enforced)

---

## Rate Limiting & Retry Logic

### Command Timeout
- HTTP command requests timeout after **30 seconds**
- WebSocket command subscriptions timeout after **30 seconds**

### Event Handling Retry
- Automatic restart on `UNEXPECTED_STOP` events
- Maximum **3 retry attempts** per machine
- Retry state expires after **10 minutes**

---

## WebSocket Subscription Model

The WebSocket service supports a pub-sub pattern where:

1. **Agents** connect and are registered by machine ID
2. **Users** can subscribe to machine events/streams
3. **Commands** are routed to specific machines
4. **Responses** are sent back to waiting subscribers
5. **Streams** (metrics) are broadcast to all subscribers

---

## Examples

### Complete Authentication Flow

```bash
# 1. Sign up
curl -X POST http://localhost:8080/api/auth/user/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret","email":"john@example.com"}'

# 2. Login
curl -X POST http://localhost:8080/api/auth/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret"}'

# Response: {"message":"Login successful","user_id":1,"token":"eyJ..."}

# 3. Use token in subsequent requests
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer eyJ..."
```

### Machine Registration & Command

```bash
# 1. Register a machine
curl -X POST http://localhost:8080/api/user/machine \
  -H "Authorization: Bearer eyJ..." \
  -H "Content-Type: application/json" \
  -d '{"name":"My Server","users":[1],"ip":"192.168.1.100"}'

# 2. List containers on the machine
curl -X POST http://localhost:8080/api/user/cmd \
  -H "Authorization: Bearer eyJ..." \
  -H "Content-Type: application/json" \
  -d '{
    "machine_id":"agent-1",
    "command_type":"list_container",
    "params":{"list":{"all":true}}
  }'
```

---

## Configuration

The API requires the following configuration (typically in environment variables or config file):

- `Port` - Server port (e.g., ":8080")
- `USER_JWT_SECRET` - Secret key for JWT token generation
- `AllowMethods` - CORS allowed methods
- `AllowHeaders` - CORS allowed headers
- `ExposeHeaders` - CORS exposed headers
- `AllowCredentials` - CORS credentials flag

---

## Notes

- Passwords are currently stored in plain text - **implement hashing in production**
- CORS is configured to allow all origins - **restrict in production**
- Machine middleware is defined but not currently enforced
- User middleware is defined but not currently enforced on user endpoints
- Database transactions should be implemented for critical operations
- Input validation can be enhanced with more specific validators

---

## Future Enhancements

- Password hashing (bcrypt)
- Email verification
- Password reset functionality
- API rate limiting
- Request logging and monitoring
- API versioning
- Pagination for list endpoints
- Advanced filtering and sorting
- Batch operations
- Webhook support for events
- Real-time notifications

---

## Support

For issues or questions, please contact the development team or file an issue in the repository.

**Last Updated:** February 1, 2026
