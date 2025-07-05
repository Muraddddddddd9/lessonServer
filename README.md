# Lesson Server

The server-side component of a system for managing lessons, tests, and time synchronization between participants.

## Features
- User login and registration  
- Lesson stage management (start, pause, switch)  
- Time control for lessons and tests  
- WebSocket implementation for real-time time synchronization  
- Test handling: retrieving questions, submitting answers  
- Data storage in MySQL  

## Technologies  
- Programming language: Go  
- Web framework: Fiber (for API and WebSocket)  
- Database: MySQL  

## Installation and Setup  

### Local Setup  

1. Clone the repository:  
```bash  
git clone <repository>  
```  

2. Database setup:  
- Install and start MySQL server  
- Create a database  
- Import the initial schema from `database/init.sql`  

3. Environment configuration:  
```bash  
cp .env.example .env  
```  
Edit the `.env` file according to your configuration  

4. Install dependencies:  
```bash  
go get .  
```  

5. Start the server:  
```bash  
go run .  
```  

The server will be available at: `http://localhost:8080`  

### Running with Docker Compose  

1. Ensure Docker and Docker Compose are installed  

2. Copy and configure the environment file:  
```bash  
cp .env.example .env  
```  
Edit `.env` if needed (default settings usually work for Docker Compose)  

3. Start the services:  
```bash  
docker-compose -f docker-compose.yml up -d  
```  