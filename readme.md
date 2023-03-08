# Welsh academy
This project is an example of REST API. 
This project is written in golang.
The ApI provides endpoint to: 

- Manage user (create, login, get user profile)
- Manage recipe (create, get, delete, add to favorite, remove favorite)
- Manage ingredient for a recipe (create, get, delete)

## Installation
This project use docker for the dev and the run environment. 
Te dev environment use a dev container wwith vscode. 
To start the project: 
- Copy the file .env.example and name it .env
- Update the .env information as described in the example file
- To start in production run docker-compose up
- To run in dev mode use vscode remote container extension and open the folder with the extension
## test the project
The project provide file welsh_academy.http which is a rest-client file. 
Use the file with the rest-client vscode extension 
