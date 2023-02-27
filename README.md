## Config example

```yaml
server:
  addr: :8080
  host: site.org
postgres:
  host: localhost
  port: 5432
  database: cc
  username: postgres
  password: PASSWORD
auth:
  expiration_time: 1h
  signing_key: SIGNING_KEY
```