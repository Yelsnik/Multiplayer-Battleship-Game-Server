# Multiplayer Battleship Game Server

A **real-time multiplayer Battleship game server** built with **Go** and **WebSockets** for low-latency networking, game state synchronization, and concurrency control.

## Features

- 🛳 **Multiplayer Gameplay** – Players can challenge each other in real-time.
- ⚡ **Low Latency** – Uses WebSockets for fast, two-way communication.
- 📡 **Game State Synchronization** – Ensures all players see the same board state.
- 🔒 **Concurrency Control** – Prevents race conditions with proper locking mechanisms.
- 🎯 **Turn-Based System** – Players take turns attacking and defending.
- 🏗 **Scalable Architecture** – Supports multiple concurrent matches.

## Tech Stack

- **Backend:** Go (Gorilla WebSocket, gRPC)
- **Networking:** WebSockets for real-time updates
- **Database (Optional):** Redis/PostgreSQL for persistence (if needed)
- **Containerization:** Docker (for easy deployment)

## Installation

1. **Clone the Repository**
   ```sh
   git clone https://github.com/Yelsnik/Multiplayer-Battleship-Game-Server.git
   cd game-sever
   ```


2. **Install Dependencies for both client and sever**
   ```sh
   go mod tidy
   ```

3. **Run the server**
    ```
    cd game-severs
    make run
    ```

4. **Run the client**
    ```
    cd game-clients
    make run
    ```




## How It Works

1. **Players Connect via WebSockets** – Clients connect to the server via `ws://localhost:8000/ws`.
2. **Matchmaking System** – Players are matched with opponents automatically.
3. **Game Board Initialization** – Each player places ships before the game starts.
4. **Turn-Based Attacks** – Players take turns attacking a grid position.
5. **Game State Updates** – The server broadcasts updates to both players.
6. **Win Condition** – The first player to sink all opponent ships wins.

<!-- ## API Endpoints -->

<!-- | Method | Endpoint      | Description                    |
| ------ | ------------- | ------------------------------ |
| `WS`   | `/ws`         | WebSocket connection endpoint  |
| `POST` | `/start-game` | Initiates a new game           |
| `POST` | `/attack`     | Player attacks a grid cell     |
| `GET`  | `/game-state` | Fetches the current game state | -->
<!-- 
## WebSocket Message Format

- **Client sends (attack request):**
  ```json
  {
    "type": "attack",
    "x": 3,
    "y": 5
  }
  ```
- **Server responds:**
  ```json
  {
    "type": "attack_result",
    "hit": true,
    "sunk": false
  }
  ``` -->

## Contributing

1. **Fork the repository**
2. **Create a feature branch**
   ```sh
   git checkout -b feature-name
   ```
3. **Commit your changes**
   ```sh
   git commit -m "Add new feature"
   ```
4. **Push to GitHub**
   ```sh
   git push origin feature-name
   ```
5. **Create a Pull Request**

## License

This project is licensed under the **MIT License**.

## Contact

For questions or suggestions, reach out to **[your email]** or open an issue in the repository.

---

🚀 **Happy Battleshipping!**

