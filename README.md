# Field Service Management

## Overview

This is my capstone project for Boot.Dev.

The goal of the project was to use as much of what I learnt from Boot.Dev.

The API server is served with Go.

The database uses Postgres.

The front-end uses HTML, CSS and JavaScript.

All of if is contained in Docker containers.

## Components of Project

### Database

I am using a dockerize Postgres 18 database for the handling of information.

### API Server

I am using a dockerized API server built in Go to handle the connection between the database and the front-end.

### Front-end

The front-end is served by the API server. 
The webpages are built in HTML and CSS.
And I used JavaScript to handle the response and request to the API Server and to handle the logic behind the webpages.