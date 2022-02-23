# Multiplayer Room Backend

A serverless back-end for managing rooms for a multiplayer lobby/room based game

## Description

This is a simple backend using Google Cloud Functions to allow users to create Rooms, join an existing Room, leave their current Room or query all available Rooms. The functions also create User objects storing the some basic user details, and each room records the users that are connected to that room. Rooms and Users are collections of Documents hosted on Firestore.

## Disclaimer
This was written as a quick one-day prototype to explore a system for publishing and discovering rooms for multiplayer games. This is my first project using Go, Protobuf, Google Cloud Functions and Firestore, so i don't gurantee any reliability or security of the code in this project. Creating projects with Firestore and Google Cloud Functions requires a Google Cloud account registered with a credit card number for identification. 

## Getting Started

### Dependencies
- Go (https://go.dev/dl/)
- Node JS (https://nodejs.org/en/download/)
- Google Cloud SDK & CLI (https://cloud.google.com/sdk/docs/install)
- Protobuf (https://github.com/protocolbuffers/protobuf/releases/latest)
- Go support for protocol buffers (https://github.com/protocolbuffers/protobuf-go/releases)

### Setup
Since this backend is created with Google Cloud Functions and Google Firestore, you'll need to setup a project on Google Cloud Platform. This will require a credit card number for proof of identification, but they give you free trial credits for individuals just prototyping.
- Install the dependencies listed above
- Create a Google Cloud Platform account or login to an existing one (https://cloud.google.com/)
- Create a new project
- Enable Firestore and Google Cloud Functions
- Run 'gcloud init' and login to your google account so that the gcloud CLI can access Google Cloud
- Ensure that GOPATH, GOROOT and GOBIN are in your PATH (https://www.programming-books.io/essential/go/gopath-goroot-gobin-d6da4b8481f94757bae43be1fdfa9e73)
- Get all Go dependencies (cd to go folder and run 'go get -d ./...')

### Deployment
- Compile the proto files (Run '.\buildproto.bat')
- Deploy the cloud functions (Run '.\deploy.bat RoomCreate', '.\deploy.bat RoomJoin', '.\deploy.bat RoomLeave', '.\deploy.bat RoomGetAll')

### Testing
You can run the Google Cloud Functions via the Google Cloud dashboard. Under the Cloud Functions section, click your functions hamburger menu and select 'Test Function', or click the function name and navigate to the Testing tab. The testing tab will allow you to run the function and pass in a json object. The models that the functions take as parameters can be seen in Rooms.proto and Users.proto. E.g, the RoomCreate method requires a RoomCreateRequest object parameter and returns a RoomCreateResponse object, etc.