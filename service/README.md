# Daytrack

Daytrack is a foss software for tracking activity during the day

## Create jwt ed25519 secrets

How to generate the keys and add them to your env variable

```bash
openssl genpkey -algorithm ed25519 -out private.pem
openssl pkey -in private.pem -pubout -out public.pem

cat private.pem | sed -e :a -e '/$/N; s/\n/\\n/; ta'
cat public.pem | sed -e :a -e '/$/N; s/\n/\\n/; ta'
```


## Create env file

```env
HTTP_HOST=localhost
HTTP_PORT=3000

OPAQUE_SALT_NONCE=
OPAQUE_CHALLENGE_RANGE=

JWT_PRIVATE_KEY=
JWT_PUBLIC_KEY=
JWT_ACCESS_TOKEN_DURATION=15m
```

## Generate swagger

```bash
swag init -o openapi
```

## How to run

```bash
go run .
```