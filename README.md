# FORUM 🐰

Forum is a web application that allows authorized users to create, comment, like, dislike and add tags to the their own posts.

## Features 👀
- Creating post
- Post tags (categories)
- Filter by tags
- Likes, dislikes on posts and on comments
- Commenting on posts
- Registration
- Authorization
- Database connection

## Requirements 🥺

1. Go 1.20 or higher
2. Have an access to this repository 🙂

## Usage 😳

**To start using this application follow these steps:**

1. Clone the repository   
```bash     
    git clone git@git.01.alem.school:msarvaro/forum.git 
```
2. Download dependecies:
```go
    go mod tidy
```
3. Run these commands in terminal (running the server in docker)
```bash
    make build
    make docker-run
```

4. By default the server would run by next address - http://localhost:5000

- If you've finished reviewing, please delete the docker image and container by running next command:
```bash
    make clean
```

---
**Alternatively** 🔄

You can run the server directly:
```bash
    make run
```
---

## Test account 🧪

Regular user:
- username - **"nah"**
- email - **"naaah@nah.com"**
- password - **"nahnahnah"**

## NOTES 📝

- You can react on post only from post page
- You CANNOT create a new user using non-ascii symbols in either username or email
- You can authorize either with **username** or **email**
- If don't have make installed, you may run the app with next command (**run from root directory**):
```go
    go run ./cmd/app
```

## Authors 👁️👅👁️

[msarvaro](https://01.alem.school/git/msarvaro)

[nmagau](https://01.alem.school/git/nmagau)