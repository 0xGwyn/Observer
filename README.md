# Observer
Observer is a simple tool that checks whether public assets of companies in hackerone, bugcrowd, intigriti and yeswehack have changed, then 
pushes notifications to the discord webhook url provided. It works based on the https://github.com/Osb0rn3/bugbounty-targets repository.

# Usage
The `discordURL_VDP` and `discordURL_BBP` variables inside the main.go file have to be given the correct discord webhook urls.
