dependencies:
	@command -v fresh --version >/dev/null 2>&1 || { printf >&2 "fresh is not installed, please run: go get github.com/pilu/fresh\n"; exit 1; }

serve: dependencies
	@sh -c "MONGO_URL=127.0.0.1 MONGO_DB_NAME=killer-koala PORT=9090 PRIVATE_KEY=koala_key.pem PUBLIC_KEY=koala_key.pub fresh -c runner.conf"
