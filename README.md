# spinner
Maintain socket connections to clients to spin up instances.

## Base Functionality
* Keep a list of machines (captain) by open socket connect
* Ferry data between captains and nebula central

## Possible Removal
The spinner may be later deprecated and removed.  
It's functionality being replaced by:  
* Captains hearbeat to a storage system
* IP as a service (like ngrok)  

In this case the classic spinner will become purely a ferry. A ferry
that is agnostic to what underlying data and components it interacts with. 
