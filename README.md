# pdf-service

This service is part of the Linkerify project

1. Service for working with links and storing them in a database as html content: [link-saver-api](https://github.com/0x0FACED/link-saver-api)
2. Bot for telegram [link-saver-telegram-bot](https://github.com/0x0FACED/link-saver-telegram-bot)


## Task

The goal of this service is to allow users to save the content of HTML pages as pdf files. 

The problem of saving html content and then returning the link to the user is that we cannot save dynamic sites properly. 

In addition there are resource issues. 

**The ultimatum solution is to save as a pdf file.**

## TODO

- [ ] Implement service methods
- [ ] Add go-rod implementation
- [ ] Add redis and organize storage
- [ ] Write a docker file
- [ ] Add tests


<h1>
  <p align="center">
<strong>Work in progress</strong><br/>
</h1>