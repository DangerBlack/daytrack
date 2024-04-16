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
