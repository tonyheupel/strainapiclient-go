# The Strain API Go Module

[![Build Status](https://travis-ci.com/tchype/strainapiclient-go.svg?branch=master)](https://travis-ci.com/tchype/strainapiclient-go)

# Background
 This is Go module that interacts with The Strain API. The site can be found at https://strains.evanbusse.com/
 which is where you can sign up for an API Key for free (at the time of this writing).

 The API base URL is at https://strainapi.evanbusse.com/API_KEY/ where `API_KEY` is the API Key you are issued.

 # Usage
 WIP

 # Additional Features

## Extensibility

## Implement your own Client

 You can build your own `Client` by simply creating a struct and functions that implement the `Client` interface.
 This module comes with its own default client, called (unimaginitively) `DefaultClient`

## Use your own handler for API requests from the DefaultClient

 If you don't want to fully implement your own `Client`, you can simply provide your own function 
 to handle the requests from the `DefaultClient`. This is especially useful for testing when you want 
 to simply provide your own custom JSON as the result of the API call.

 You can see [a couple of examples](./examples/default_client/simple.go):
 1. where we override all requests to return an error whose message includes the path that was requested; good for unit testing and tracking calls
 1. where we only override one requst for a specific operation on a specific strain ID, but leave the rest of the existing logic and api calls intact.
 1. more to come...
