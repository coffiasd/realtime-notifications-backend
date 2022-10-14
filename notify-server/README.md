# notifty-server

## monitor

Monitor blockchain realtime transaction via algorand REST API.
Starts from current round.
And then sleep 1 second wait for next request with parameter round increase 1.
Fetch the response transations json data.
If the current transaction match events that our user subscribe,let's say the transaction receiver is our subscription user.
That means his/her transaction is compelete,our server will push a notification message to his/her mobile client.
At the same time if our subscribers offered an email,our server will send a mail too.
