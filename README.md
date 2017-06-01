Chrom(e|ium) Screens
===

Chome/Chromium-based dashboards, changeable via Slackbots.

# Usage

Copy on to target PC. First run only requires a Slack API token env var - this should get saved to `config.json` in the same directory.

```
$ SLACK_API_TOKEN=xoxs-... ./chrome-screens
```

Once it is connected it should tell you the name of the user it connected as (configured via Slack bot configuration). Invite that user to any channel you want to issue commands via.

Once the Slack bot is in a channel you can ask it for help: `@mybotname help`. You should be able to work it out from there.

# Installing

If you have Go installed:

```
go get -u github.com/porty/chrome-screens
```

If you don't have Go installed: install it on another machine, scp it to the target. Cross compile if needed.

# Getting a Slack token

Slack API tokens can be retrieved by hitting _Apps & integrations_, search for _Bots_ (should be the first hit, with a boring grey icon) and configuring a new bot. This will give you a Bot integration with an API token. Be sure to give it a good icon/avatar and a name.

# Windows support

I'm lazy. I don't even know if this supports OSX.

# Licence

MIT
