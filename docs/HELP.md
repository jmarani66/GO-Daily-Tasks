# Daily Task Planner - Help

This document provides step-by-step instructions on how to run the Daily Task Planner application on a Debian Linux desktop.

## Prerequisites

- Go (version 1.15 or later)
- Git

## Installation

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd daily-task-planner
   ```

2. **Build and run the backend:**

   Open a terminal and navigate to the `backend` directory:

   ```bash
   cd backend
   ```

   Run the application:

   ```bash
   go run main.go
   ```

   The server will start on port 8080.

## Accessing the Application

Open your web browser and go to the following URL:

[http://localhost:8080](http://localhost:8080)

You should see the Daily Task Planner application in your browser.

## How to Use

- **Add a task:** Type a task in the input field and click the "Add" button or press Enter.
- **Mark a task as complete:** Click the checkbox next to a task.
- **Delete a task:** Click the "Delete" button next to a task.
