# FORUM 🐰 (1.0)

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
- Image upload (.jpg, .png, .gif, .jpeg)
- Authentication through Google or Github
- Notification system on most of user actions (like/dislike post or comment, comment post, report comment, delete comment, request for moderator role status, declining requests, etc.)
- Your reacted/commented/created posts pages
- Admin panel - reports, role upgrade requests, all users
- Tag creation, deletion (only by admin)
- Report post or comment (moderator role required)

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

4. By default the server would run by next address - https://localhost:5000

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

Admin:
- username - **nah**
- email - **naaah@nah.com**
- password - **nahnahnah**

## NOTES 📝

- You can react on post only from post page
- You CANNOT create a new user using non-ascii symbols in either username or email
- You can authorize either with **username** or **email**
- Admin can respond to moderator's report only by accepting it (deleting reported content) or by declining it (notify report's author that it has been declined)
- If don't have make installed, you may run the app with next command (**run from root directory**):
```go
    go run ./cmd/app
```

## Authors 👁️👅👁️

[msarvaro](https://01.alem.school/git/msarvaro)

[nmagau](https://01.alem.school/git/nmagau)

<span style="color: #999;">huge update 2.0 is coming soon...</span>

