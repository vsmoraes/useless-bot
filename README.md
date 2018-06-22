# Useless bot

This is just another simple Slack bot

## What was done so far?
* Simple `hey` command: The bot will send a simple response to you;
* `tweet` command: if your keys are correct, the bot can send a tweet using this command;


## `hey` command
This is almost like a "hello world" command. With the bot properly configured, you can send it a message and it'll reply, like this:

```
<user> @useless_bot hey you!
<userless_bot> @user: hey there!
```

## `tweet` command
You can send a tweet to the configured account using this command. Just:

```
<user> @useless_bot tweet This is a testing tweet...
<useless_bot> @user: tweet sent!
              This is a testing tweet...
```

## Commands not recognized
There's a simple command check to validate the available commands, if a not recognized command reaches the bot a response will be dispatched like this:

```
<user> @useless_bot post this on twitter: foo
<useless_bot> @user: command not recognized: `post`
```

As you can see, we just validate this second word (separeted by empty spaces) of your message. So expect many bugs if not used just like the docs are showing ;-)

# TODOS
- [ ] increase the test coverage
- [ ] send a failed tweet attempt to slack
- [ ] better validation of commands (e.g.: you can mention the bot's user at the end of your message and it still works, which sucks)
- [ ] add new commands?
- [ ] make possible to extend the command list and handler as easy as possible
- [ ] write a full integration and configuration with Slack
- [ ] write a how-to on how to deploy this project on Heroku so you can use in production
