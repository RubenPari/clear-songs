# Clear Songs

Clear Songs is a REST API service that helps you efficiently manage your Spotify music library by providing powerful bulk deletion capabilities for your tracks and playlists.

## Overview

The application allows you to quickly clean up your Spotify library by:
- Removing all tracks from a specific artist
- Deleting tracks based on the number of songs per artist
- Clearing out entire playlists while maintaining the playlist structure

## Features

### Track Management
- **Delete by Artist**: Remove all tracks from a specific artist in your library
- **Quantitative Deletion**: Delete tracks based on the number of songs you have per artist (e.g., remove all tracks from artists with more than X songs)
- **Backup System**: Automatically saves deleted tracks to MySQL database for recovery in case of accidental deletion

### Playlist Management
- **Playlist Clearing**: Empty any playlist you own while keeping the playlist itself intact
- **Bulk Operations**: Perform operations quickly and efficiently through the API

## Setup

### Environment Variables
Create a `.env` file in the root directory with the following parameters:

```env
# Spotify API Credentials
CLIENT_ID=your_spotify_client_id
CLIENT_SECRET=your_spotify_client_secret
REDIRECT_URL=your_callback_url

# Database Configuration
DB_HOST=your_database_host
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
DB_PORT=your_database_port
```

### Database Setup
The application uses MySQL as its database system for track backup functionality. This feature allows you to:
- Automatically save tracks before deletion
- Recover accidentally deleted tracks
- Maintain a history of your library changes

You can use any MySQL database of your choice. The application will create the necessary tables on startup.

## API Endpoints

### Authentication
- `/auth/login` - Initiates Spotify OAuth flow
- `/auth/callback` - Handles OAuth callback from Spotify
- `/auth/logout` - Logs out the current user
- `/auth/status` - Checks authentication status

### Track Operations
- `/track/artist/{id_artist}` - Delete all tracks from a specific artist
- `/track/range` - Delete tracks based on quantity parameters

### Playlist Operations
- `/playlist/tracks` - Remove all tracks from a specified playlist
- `/playlist/tracks/all` - Remove tracks from both playlist and user library

### Album Operations
- `/album/convert` - Convert an album to individual songs in your library

## Getting Started

1. Clone the repository
2. Create and configure your `.env` file with all required parameters
3. Set up your MySQL database
4. Configure your Spotify API credentials in the Spotify Developer Dashboard
5. Start the server
6. Authenticate with your Spotify account
7. Begin managing your library through the API endpoints

## Security

The application uses OAuth 2.0 for authentication with Spotify and requires appropriate permissions to manage your library and playlists.

## Technical Stack

- Go (Golang)
- Gin Web Framework
- Spotify Web API
- Redis for caching
- MySQL for track backup storage

**Note**: While the application includes a backup system, it's still recommended to be careful when using bulk deletion features. Always verify your selections before performing bulk operations.
