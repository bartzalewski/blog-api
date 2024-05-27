# Blog API

Welcome to the Blog API! This project is a simple backend for a blogging platform, allowing user management, post creation, and commenting features.

## Features

- User authentication with JWT tokens
- CRUD operations for blog posts
- Adding comments to posts
- Secure password storage using bcrypt
- Concurrent safe in-memory storage for users, posts, and comments

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed on your system:

- [Go](https://golang.org/doc/install) (version 1.16+)
- [Git](https://git-scm.com/)

### Installation

1. **Clone the repository**:

   ```sh
   git clone https://github.com/bartzalewski/blog-api.git
   cd blog-api
   ```

2. **Initialize Go modules**:

   ```sh
   go mod tidy
   ```

3. **Run the application**:

   ```sh
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## API Endpoints

### User Authentication

- **Sign Up**

  `POST /signup`

  Request:

  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```

  Response:

  ```json
  {
    "status": "User created successfully"
  }
  ```

- **Sign In**

  `POST /signin`

  Request:

  ```json
  {
    "username": "your-username",
    "password": "your-password"
  }
  ```

  Response:

  ```json
  {
    "status": "Signed in successfully"
  }
  ```

  Cookie: `token`

### Blog Management

- **Create Post**

  `POST /posts`

  Request:

  ```json
  {
    "title": "First Post",
    "content": "This is the content of the first post."
  }
  ```

  Response:

  ```json
  {
    "id": 1,
    "title": "First Post",
    "content": "This is the content of the first post.",
    "author": "your-username",
    "created_at": "2023-05-24T12:34:56Z",
    "comments": []
  }
  ```

- **Get Posts**

  `GET /posts`

  Response:

  ```json
  [
    {
      "id": 1,
      "title": "First Post",
      "content": "This is the content of the first post.",
      "author": "your-username",
      "created_at": "2023-05-24T12:34:56Z",
      "comments": []
    }
  ]
  ```

- **Add Comment**

  `POST /posts/{id}/comments`

  Request:

  ```json
  {
    "content": "This is a comment."
  }
  ```

  Response:

  ```json
  {
    "id": 1,
    "content": "This is a comment.",
    "author": "your-username",
    "created_at": "2023-05-24T12:34:56Z"
  }
  ```

## Built With

- [Go](https://golang.org/) - The Go programming language
- [Gorilla Mux](https://github.com/gorilla/mux) - A powerful URL router and dispatcher for Golang
- [JWT-Go](https://github.com/dgrijalva/jwt-go) - A Go implementation of JSON Web Tokens
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) - A package for password hashing

## Contributing

Feel free to submit issues or pull requests. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the Go community for their invaluable resources and support.
- [Gorilla Mux](https://github.com/gorilla/mux) and [JWT-Go](https://github.com/dgrijalva/jwt-go) for making development easier.

---

Happy coding! 🚀
