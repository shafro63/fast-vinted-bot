
# Fast Vinted Monitor 

Never miss an item on vinted with this simple and fast monitor ! 🚀



## Features

* Be the first to see the latest items posted on vinted 
* Items sent via discord  
* Custom links supported
* Unlimited monitoring links 
* Saved sessions and user's data in database
* All countries supported

## Requirements

* Your own computer or VPS (Just in case you want 24/24 run)
* A discord server
* [A discord Bot](https://discord.com/developers/applications) (Guild install only with admin permission)
* [MongoDB Community installed](https://www.mongodb.com/docs/manual/administration/install-community/#std-label-install-mdb-community-edition)
* [Golang installed](https://go.dev/doc/install)
* Git installed (optionnal)
* **A good proxy provider** for rotating proxies — I recommend using Decodo(ex Smartproxy) 99% Success rate. Theirs datacenter proxies work like a charm and are their cheapest plan. Pay less with [my link](dashboard.decodo.com/register?referral_code=cbc3f26dd85fd04cfef11103b48e1360eef6b23f)

## Setup

We assume you already have the requirements — If not please don't go further before.

1. Invite your discord bot to your server if not done yet
2. Clone the repository or download the code
3. Launch mongoDB in the background (You can see how to run that for your device in the installation webpage —You can even set to run it automatically at startup so you won't have to to do it every single time later)
4. Open the .env file and put the **mandatory** credentials 
5. To launch the monitor :
Go to the app folder directory from the terminal

example: cd /Users/login/Desktop/fast-vinted-bot 
```bash
  cd <folderpath> 
```
Then
```bash
  go run . 
```



## Features

* Be the first to see the latest items posted on vinted 
* Items sent via discord  
* Custom links supported
* Unlimited monitoring links 
* Saved sessions and user's data in database
* All countries supported

## Usage
Open Vinted, select your filters and copy the link address

Your vinted link must follow the path /catalog?

Once the discord bot added to your server and the bot is launched you will see the commands names when tapping on " / "

* /create_private_channel : Create a private channel 
* /add_link : Add your vinted link to one of your channel — This will instantly start monitoring the link 
* /delete_private_channel : Delete private channel
* /delete_link : Delete one link from your private channel

To visually see the server sessions & data use [mongoDB Compass](https://www.mongodb.com/try/download/compass)
## Note

The monitor refresh the link every 2 seconds, you can change that in your env file

If you want more features don't hesitate to open an issue :)

I will add code's comments later. For now blame my lazyness :(