# Go Chat Application

This is a chat application built with Go, featuring both a WebSocket-based web interface and a TCP-based console client/server. It also includes a separate user authentication module.

## Features

  * **Real-time Chat:**
      * WebSocket support for browser-based clients.
      * TCP socket support for console-based clients.
  * **User Authentication (Separate Module):**
      * User signup
      * User login
      * MySQL database for user storage

## Project Structure

The project is organized into two main components:

1.  `socket_chat`: Contains the core chat functionalities.
      * `client/`: TCP console client.
      * `server/`: TCP server.
      * `web/`: WebSocket server and HTML/JavaScript client.
2.  `login_go`: Contains the user authentication service.
      * `dao/`: Data Access Object for interacting with the database.

## Prerequisites

  * Go (version 1.x or higher)
  * MySQL (for the `login_go` module)

## Setup and Running

### 1\. Socket Chat

#### a. WebSocket (Web) Chat

This version allows users to chat through a web browser.

**Running the Web Server:**

1.  Navigate to the `socket_chat/web` directory:
    ```bash
    cd go-backend-main/socket_chat/web
    ```
2.  Run the Go application:
    ```bash
    go run main.go
    ```
3.  Open `websockets.html` (or navigate to `http://localhost:8080/`) in your web browser. The server will echo back any message you send.

#### b. TCP Console Chat

This version allows users to chat through a console/terminal.

**Running the TCP Server:**

1.  Navigate to the `socket_chat/server` directory:
    ```bash
    cd go-backend-main/socket_chat/server
    ```
2.  Run the Go application:
    ```bash
    go run main.go
    ```
    The server will start listening on `localhost:8080`.

**Running the TCP Client:**

1.  Open a new terminal.
2.  Navigate to the `socket_chat/client` directory:
    ```bash
    cd go-backend-main/socket_chat/client
    ```
3.  Run the Go application:
    ```bash
    go run main.go
    ```
4.  Enter your username when prompted.
5.  You can now send and receive messages. The client supports reconnection attempts and heartbeats.

### 2\. Login System (`login_go`)

This module provides user signup and login functionality. It interacts with a MySQL database.

**Database Setup:**

1.  Ensure MySQL is installed and running.
2.  Update the database connection string in `login_go/dao/dao.go` if necessary:
    ```go
    db, err := sql.Open("mysql", "root:123456@(127.0.0.1:3306)/pwdatabase?parseTime=true")
    ```
    Modify `root:123456`, `127.0.0.1:3306`, and `pwdatabase` according to your MySQL setup.
3.  The application will automatically create the `users` table if it doesn't exist.

**Running the Login Server:**

1.  Navigate to the `login_go` directory:
    ```bash
    cd go-backend-main/login_go
    ```
2.  Run the Go application:
    ```bash
    go run main.go
    ```
    The server will start listening on port `:8080` (and also attempts `:80`).

**API Endpoints:**

  * **Sign Up:** `POST /user/api/signup`
      * Request Body (JSON): `{"username": "your_username", "password": "your_password"}`
  * **Log In:** `POST /user/api/login`
      * Request Body (JSON): `{"username": "your_username", "password": "your_password"}`

## Dependencies

The application uses the following Go packages:

  * **`socket_chat/web`:**
      * `github.com/gorilla/websocket`
  * **`login_go`:**
      * `github.com/google/uuid`
      * `github.com/gorilla/mux`
      * `github.com/gorilla/sessions`
      * `golang.org/x/crypto/bcrypt`
      * `github.com/go-sql-driver/mysql` (for the DAO)

You can install these dependencies using `go get`. For example:

```bash
go get github.com/gorilla/websocket
go get github.com/google/uuid
# and so on for other dependencies
```

## How it Works

### Socket Chat (Web)

  * The server (`socket_chat/web/main.go`) uses the `gorilla/websocket` package to upgrade HTTP connections to WebSocket connections on the `/echo` endpoint.
  * It serves an HTML file (`websockets.html`) which contains JavaScript to establish a WebSocket connection to `ws://localhost:8080/echo`.
  * Messages sent from the client are received by the server and echoed back to the same client.

### Socket Chat (TCP)

  * **Server (`socket_chat/server/main.go`):**
      * Listens for incoming TCP connections on `localhost:8080`.
      * Manages a map of connected clients.
      * When a message is received from a client (JSON format: `{"username": "user", "msgType": "chat", "msgContent": "hello"}`), it's decoded.
      * If the `msgType` is "chat", the message is broadcast to all connected clients.
  * **Client (`socket_chat/client/main.go`):**
      * Prompts the user for a username.
      * Connects to the TCP server at `localhost:8080`.
      * Handles reading messages from the server in one goroutine and sending user input (as JSON messages) in another.
      * Implements a heartbeat mechanism to keep the connection alive.
      * Includes a retry mechanism for reconnections if the server connection is lost.

### Login System (`login_go`)

  * Uses `gorilla/mux` for routing HTTP requests.
  * `/user/api/signup`:
      * Accepts a username and password (JSON).
      * Checks if the username already exists in the database (`dao.CheckUserExist`).
      * If not, it creates a new user (`dao.CreateUser`). The password is currently stored as plain text (Note: `HashPassword` and `CheckPasswordHash` functions exist but are marked as "not used").
  * `/user/api/login`:
      * Accepts a username and password (JSON).
      * Validates credentials against the database (`dao.CheckUserPassword`).
  * **DAO (`login_go/dao/dao.go`):**
      * Manages the connection to a MySQL database.
      * Provides functions to create the user table, create a user, check if a user exists, and validate user passwords.
