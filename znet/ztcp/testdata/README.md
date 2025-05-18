# About test data

- key.pem
- cert.pem

```bash
openssl req -new -x509 -nodes -days 36500 -subj '/CN=test' -keyout key.pem -out cert.pem
```
