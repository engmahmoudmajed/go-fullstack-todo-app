Guide to use rest exetantion

# 🚀 API Testing with REST Client

This project uses the **REST Client** extension for Visual Studio Code to test API endpoints. This allows us to keep our API documentation and testing right alongside our source code.

## 1. Installation
1. Open **VS Code**.
2. Go to the **Extensions** view (Ctrl+Shift+X).
3. Search for `REST Client` (by Huachao Mao).
4. Click **Install**.

## 2. How to use
I have created a file named `api_tests.http` in the root of this project. 

1. Open `api_tests.http`.
2. Ensure your Go server is running (`air` or `go run main.go`).
3. Click the small **"Send Request"** text that appears above each HTTP method (GET, POST, etc.).
4. The response will appear in a new pane on the right.

## 3. Example Test File (`api_tests.http`)
Create a file named `api_tests.http` and paste the following:

```http
### Variables
@baseUrl = http://localhost:4000

### Get Welcome Message
GET {{baseUrl}}/
Accept: application/json

### Get All Todos
GET {{baseUrl}}/api/todos
Accept: application/json

### Create a New Todo (Example Placeholder)
POST {{baseUrl}}/api/todos
Content-Type: application/json

{
    "title": "Finish Go Backend",
    "completed": false
}