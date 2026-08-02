# Build a RESTful CRUD API with Golang Gin

Solution for the [Bloggin Platform API](https://roadmap.sh/projects/blogging-platform-api) challenge from [roadmap.sh](https://roadmap.sh/).

### Build from Source

   ```bash
   git clone https://github.com/satyamkanungo-dev/blog-reset-api.git
   cd blog-rest-api
   ```

## Development Setup

### Using Air for Live Reload

This project uses Air for live reload during development. Here's how to get started:

1. Install Air:

```bash
go install github.com/cosmtrek/air@latest
```

2. Run the application with Air:

```bash
air
```

Air will automatically:

- Watch your Go files for changes
- Rebuild your application when files change
- Restart your server automatically
- Show build errors in a colorized format

To stop Air, press `Ctrl+C` in your terminal.

The project is already configured with a `.air.toml` file that:

- Watches the `cmd/api` directory
- Excludes test files and common directories
- Includes Go files, templates, and HTML files
- Uses colorized output for better visibility

### Environment Variables

Create a `.env` file with these variables (all have sensible defaults):

```bash
DATABASE_URL = "postgres://postgres:your_password@localhost:5432/your_db_name?sslmode=disable"
PORT = 8000
JWT_SECRET=your-secret-key
```

For production, make sure to set these values through your deployment platform's environment configuration.

### Running Without Air

If you prefer not to use Air, you can run the application directly with Go:

```bash
go run ./cmd/api
```

This will start the server on `http://localhost:8000`. Note that you'll need to manually restart the server when you make changes to the code.




### API Documentation

The API documentation is available via Swagger UI at:

```
http://localhost:8000/docs/index.html
```

This interactive documentation provides:

- Complete API endpoint listing
- Request/response schemas
- Try-it-out functionality
- Authentication details


