# Awesome Project

Welcome to **Awesome Project**, a Go API built using the Fiber framework and GORM for database interaction. This project features JWT middleware for authentication and role-based authorization. Users can perform operations like user registration, login, logout, token refresh, as well as CRUD (Create, Read, Update, Delete) operations.

## Features

- User Registration: Users can create new accounts with a unique username and password.
- User Login: Registered users can log in and receive an access token and a refresh token.
- Token Refresh: Users can refresh their access token using a refresh token.
- User Logout: Users can log out and invalidate their tokens.
- Authentication: JWT middleware ensures secure authentication.
- Authorization: Role-based authorization restricts access to certain routes based on user roles.
- CRUD Operations: Perform Create, Read, Update, and Delete operations on resources.
- Fiber: Utilizes the efficient and fast Fiber web framework.
- GORM: Interacts with the database using the GORM ORM.

## Installation

1. Clone the repository:

    ```sh
    git clone [https://github.com/your-username/awesome-project.git](https://github.com/genss333/awesomeProject)
    cd awesome-project
    ```

2. Install dependencies using your preferred package manager (e.g., Go modules):

    ```sh
    go mod tidy
    ```

3. Configure the database connection in the `.env` file.

4. Run the application:

    ```sh
    go run main.go
    ```

5. Access the API at `http://localhost:3000`.

## API Endpoints

- `POST /register`: Register a new user.
- `POST /login`: Authenticate and obtain tokens.
- `POST /refresh-token`: Refresh the access token.
- `POST /logout`: Invalidate tokens and log out.
- `GET /user`: Access user resources (requires "user" role).
- `GET /resource/:id`: Get resource by ID (requires authentication).
- `POST /resource`: Create a new resource (requires "admin" role).
- `PUT /resource/`: Update resource by ID (requires "admin" role).
- `DELETE /resource/:id`: Delete resource by ID (requires "admin" role).

## Contributing

Contributions are welcome! If you find any issues or have suggestions, feel free to create a pull request or submit an issue.

Feel free to customize this README to fit your project's specific details. Make sure to replace placeholders like `your-username` and `awesome-project` with the appropriate values. Happy coding on your "Awesome Project"!

