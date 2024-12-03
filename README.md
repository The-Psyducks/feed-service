# Feed Service

[![Coverage Status](https://coveralls.io/repos/github/The-Psyducks/feed-service/badge.svg?branch=main)](https://coveralls.io/github/The-Psyducks/feed-service?branch=main)

## Table of Contents

1. [Introduction](#introduction)
2. [Pre-Requisites](#pre-requisites)
3. [How To Run](#how-to-run)
4. [Tests](#tests)

## Introduction

This repository manages all interactions with Twitsnaps, including creation, search, or presentation in various feed formats. It also supports engaging with Twitsnaps through likes, retweets and bookmarks.

## Pre-Requisites

To set up the microseervice development environment, complete the `.env` template in `server/`.

Proper functionality requires both the User-Service and the Messages-Service microservices are up and running. 

### Docker Requirements

The project runs entirely in Docker, so to start the environment, you will need:

Docker Engine: Minimum recommended version 19.x Docker Compose: Minimum recommended version 1.27 

### Local Requirements

To set up and to run the project lovally without docker, the following are necessary:

Go: Version 1.23.0
This project uses MongoDB so a local connection to a database is also required.

## How To Run

Both server and database are dockerized, so to run the project only the following commands are needed: 

    docker compose build
    docker compose up service

## Tests

Tests are also dockerized, so to run them all the commands needed are"

    docker compose build
    docker compose run test
