# reactbot
### Adds reactions to messages based on content.

The bot has one command:

    !react <on | off | reload>   Enables/disables reactions or reloads the configuration file.


Config file `reactbot.cfg` example:

    react = true # activate the bot

    [[search]]
    text = "(?i)dead\\b" # A regex match.
    reaction = ["👻"]
    regex = true         # Required if this is a regex match. Otherwise can be set to false or ommitted.

    [[search]]
    channels = ["<channelID>"] # Optional. Matches all channels if not provided.
    text = "cheeto"
    reaction = ["🦧", "🧀", "🍠" ]
    regex = false

    [[search]]
    text = "<@!123456789012345678>" # A User mention reaction.
    reaction = ["🧜‍♀️", "🛵"]

---