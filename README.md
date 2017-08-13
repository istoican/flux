# Flux

Flux is a distributed key-value store with realtime notification capabilities.

It use [github.com/hashicorp/memberlist](github.com/hashicorp/memberlist) for cluster discovery and a simple hash ring implementation for key distribution.
	
It also allows setting up watchers to monitor for key changes. For that it use the HTML content streaming.


## API Reference

### Reading a key

```bash

curl http://127.0.0.1:8000/key1?watch=true

```

### Watching key changes

```bash

curl http://127.0.0.1:8000/key1?watch=true

```
### Writing a key

```bash

curl -X POST http://127.0.0.1:8000/key1 -H "Content-Type: text/plain" -d 'value1' 

```


## License

This project is licensed under the BSD License - see the [LICENSE](LICENSE) file for details.

