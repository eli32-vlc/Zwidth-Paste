# Zwidth Paste

A simple, fast markdown-based text sharing platform built with Go.

## Features

- **Markdown Editor**: Write content using markdown with real-time preview
- **Custom URLs**: Choose your own URL or use auto-generated ones
- **Edit Codes**: Protect your entries with secure edit codes
- **Modify Codes**: Share limited access codes for content-only editing
- **Spam Protection**: Hashcash proof-of-work prevents abuse
- **Admin Dashboard**: Manage entries, search, and view statistics
- **No Registration Required**: Create entries without signing up (registration optional)

## Tech Stack

- **Backend**: Go
- **Frontend**: HTML, CSS (Lightning Design System 2 inspired), JavaScript
- **Database**: SQLite
- **Security**: Hashcash for spam prevention, bcrypt for passwords

## Installation

### Prerequisites

- Go 1.21 or higher
- SQLite3

### Setup

1. Clone the repository:
```bash
git clone https://github.com/eli32-vlc/Zwidth-Paste.git
cd Zwidth-Paste
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Edit `.env` and set your configuration:
```bash
ADMIN_USERNAME=your_admin_username
ADMIN_PASSWORD=your_secure_password
SESSION_SECRET=your_random_session_secret
PORT=8080
HOST=localhost
DB_PATH=./zwidth.db
```

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
go run main.go
```

6. Open your browser and navigate to `http://localhost:8080`

## Usage

### Creating an Entry

1. Visit the homepage
2. Enter your markdown content
3. (Optional) Set a custom URL
4. (Optional) Set an edit code (otherwise one will be generated)
5. Wait for hashcash computation (prevents spam)
6. Click "Create Entry"
7. Save your edit code!

### Editing an Entry

1. Visit `/edit/{url}`
2. Enter your edit code or modify code
3. Make your changes
4. Click "Save Changes"

### Admin Dashboard

1. Visit `/admin/login`
2. Enter admin credentials (from `.env`)
3. Access dashboard at `/admin`
4. Search entries, view statistics, delete entries

## API Endpoints

- `GET /` - Homepage (create new entry)
- `POST /create` - Create new entry
- `GET /{url}` - View entry
- `GET /edit/{url}` - Edit entry form
- `POST /edit/{url}` - Update entry
- `POST /delete/{url}` - Delete entry
- `GET /raw/{url}` - Get raw markdown
- `GET /hashcash` - Get hashcash challenge
- `GET /register` - Registration form
- `POST /register` - Register user
- `GET /login` - Login form
- `POST /login` - Login user
- `GET /admin/login` - Admin login
- `GET /admin` - Admin dashboard
- `GET /admin/search` - Search entries
- `POST /admin/delete` - Delete entry (admin)

## Configuration

### Environment Variables

- `ADMIN_USERNAME`: Admin username for dashboard access
- `ADMIN_PASSWORD`: Admin password for dashboard access
- `SESSION_SECRET`: Secret key for session encryption
- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: localhost)
- `DB_PATH`: Path to SQLite database file (default: ./zwidth.db)

### Content Limits

- Text content: 200,000 characters
- Custom URL: 2-100 characters (lowercase letters, numbers, hyphens, underscores)
- Edit code: 1-100 characters
- Modify code: 1-100 characters

## Security

- **Hashcash**: Proof-of-work system prevents spam on entry creation, registration, and login
- **Edit Codes**: Entries are protected by edit codes (like passwords)
- **Modify Codes**: Limited access codes for content-only editing
- **Password Hashing**: User passwords hashed with bcrypt
- **Admin Authentication**: Separate admin authentication for dashboard access

## Development

### Project Structure

```
Zwidth-Paste/
├── database/          # Database layer
├── handlers/          # HTTP handlers
├── hashcash/          # Hashcash implementation
├── markdown/          # Markdown rendering
├── utils/            # Utility functions
├── templates/        # HTML templates
├── static/           # Static assets (CSS, JS)
├── main.go           # Application entry point
├── go.mod            # Go dependencies
└── .env              # Configuration
```

### Building

```bash
go build -o zwidth-paste main.go
```

### Running Tests

```bash
go test ./...
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- Inspired by rentry.co
- Lightning Design System 2 for UI inspiration
- Goldmark for markdown rendering
