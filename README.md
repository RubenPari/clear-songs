<p align="center">
    <em>Clear songs, a tools for delete some songs from your Spotify library</em>
</p>

---

Clear songs is a REST-API project that use Spotify public api for delete traks in your library by specific option (es. by artist, by number, by gendre ecc...).

---

***NOTE***

Please make attention when use it because you can delete all tracks saved in spotify account

## Requirements

Golang 1.18

## Installation

<div class="termy">

clone the project

```console
git clone https://github.com/RubenPari/clear-songs.git
```

</div>

<div class="termy">

install dependencies

```console
go mod download
```

</div>

<div class="termy">
create a .env file in the root of the project

```console
touch .env
```

</div>

<div class="termy">

add the following variables in the .env file

```console
CLIENT_ID=your-spotify-client-id
CLIENT_SECRET=your-spotify-client-secret
REDIRECT_URL=your-spotify-redirect-uri
PORT=your-server-port
```

</div>

## Run

<div class="termy">

```console
cd src
go run .
```

</div>

## Usage

<div class="termy">

open your browser and login to

```console
http://localhost:{PORT}/auth/login
```

explore the api with openapi to

```console
http://localhost:{PORT}/openapi/index.html
```

</div>
