# reactbot
### Adds reactions to messages based on content.

The bot has one command:

    !react <on | off>   enables/disables reactions


Config file `reactbot.cfg` example:

    react = true # activate the bot

    [[search]]
    text = "(?i)dead\\b" # regex match
    reaction = ["👻"]
    regex = true

    [[search]]
    text = "cheeto"
    reaction = ["🦧", "🧀", "🍠" ]
    regex = false

    [[search]]
    text = "<@!123456789012345678>" # user mention react
    reaction = ["🧜‍♀️", "🛵"]
    regex = false
