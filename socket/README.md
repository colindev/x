# Socket Stream Adaptor

socket 文字流轉接器

### Build

```sh
go build -o socket
```

### Usage

1. Handle script and serve for other adaptor

```sh
./socket -out 127.0.0.1:8000 /your/script
```

2. Liten server and adapt to script

```sh
./socket -in 127.0.0.1:8000 /your/script
```

### Example

1. create a file name `time-server.sh`, then `chmod +x time-server.sh`

```sh
#!/bin/bash

while true; do
    date
    sleep 5
done
```

2. and run socket server adaptor

```sh
./socket -out 127.0.0.1:8000 ./time-server.sh
```

3. run socket client adaptor

create a file name `display.php`

```php
<?php

echo fgets(STDIN);
```

```sh
./socket -in 127.0.0.1:8000 php ./display.php
```

or just run `./socket -in 127.0.0.1 cat -`

