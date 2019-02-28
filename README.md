# xchang-reminder
Send desktop notifications for upcoming calendar events.

## Configuration
Configuration is done using a yaml file and the ```-config``` command-line option.  The NTLM credentials are identified with ```domain```, ```username``` and ```password```.  ```maxfetchsize``` determines the maximum number of Calendar events to query for, ```exchangeurl``` is the EWS endpoint on your Excahnge server. ```reminder``` determines how many minutes before an event start time you will get a desktop notification for.

```yaml
---
domain: oz
username: big.kahuna
password: charge!
maxfetchsize: 5
exchangeurl: https://mail.oz.com/EWS/Exchange.asmx
reminder: 10
```

You will need to setup a crontab entry to execute xchang-reminder at regular intervals.