# A sars-cov-2 tracker bot for discord



| Command  | Description  | Req. Op  |
|---|---|---|
| !cov [country \| all]  | report the latest stats, defaults to 'usa'  | no  |
| !reaper [channel id \| off] | periodically report the death toll to the channel given or currren channel  | yes  |
| !covchans | print the channels in which the bot operates via direct message | yes  |

Config file `covidbot.cfg` example:

    TOKEN = "<rapidapitokenhere>"                        # your rapidapi token
    CHANS = ["000000000000000000", "123456789123456789"] # IDs of channels the bot will opertate in
    CRONSPEC = "1 * * * *"                               # cronspec for periodic reports
