# AVOMAC API

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Prerequisites

- GO v1.10.4
- `chi` Routing Package
- `dep` Dependancy Manager
- `MySQL` Database
- `dotenv` dotenv

## Installing

- Copy files to your local machine
- Navigate to the file path in your computer
- You also need to have the following GO packages installed

- Using the existing .env.example template, configure your environment variables and save the file .env

## Launching the program

- Run a `go run *.go`  command in your terminal and open `:port`

## API End Points

### View All Company

    Path : /company
    Method: GET
   
    Response: 201

### Restricted company

    Path : /restricted
    Method: GET
   
    Response: 201

### Search company by ID

    Path : /company/{ID}
    Method: GET
   
    Response: 201

## Deployment

Not yet deployed to a live system

## Built With

- GoLang - The core backend language

## Authors

Kennedy Kinoti

## Acknowledgments

Hat tip to anyone whose code was used

As Winston Churchill said:
> To improve is to change, to perfect is to change often.

gcloud sql connect pwc-db-instance --user=pwc-user --quiet  RhH8J4B5CwChAyta4sC7
