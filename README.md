# Timeular cli
A simple command line tool to communicate with Zei° time tracking device. Currently I have only an integration with [Harvest](https://www.getharvest.com).

## Configuration
### Harvest
Copy the `.env.example` file to `.env` and fill the environment variables with your Harvest Oauth client parameters.

Create the file `harvest.yml` in the config directory with your harvest project and task ids
#### Example `harvest.yml`
    sides:
        1:
            project_id: "14788553"
            task_id: "8341549"
        2:
            project_id: "14788553"
            task_id: "8341551"
## Usage
    sudo go run main.go

### Harvest
- Go to http://localhost:8080/harvest/login
- Click Approve in the approval screen
    
