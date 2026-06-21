# Field Service Management

## Overview

This is my capstone project for Boot.Dev.

The goal of the project was to use as much of what I learnt from Boot.Dev.

## Components of Project

### Database

The database is built on Postgres 18.

### API Server

The API server is built in Go. It handles the connection between the database and the front-end.

### Front-end

The front-end is served by the API server. 
The webpages are built in HTML and CSS.
And I used JavaScript to handle the response and request to the API Server and to handle the logic behind the webpages.

### Docker

Both the API Server and Postgres Database are in docker container for easy installation and running.
I am using a compose file to make the creation of both container easily done.