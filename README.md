# Gets Ngrok public URL

Use with docker. Wait til public url is available and update telegram update hook with it. 

```
#!/usr/bin/env bash

NGROK_PUBLIC_URL=`ngrok-url --api-host-cmd="docker-compose port ngrok 4040"` && docker-compose exec -T app php set-webhook "$NGROK_PUBLIC_URL" &

docker-compose up
```

LICENSE MIT