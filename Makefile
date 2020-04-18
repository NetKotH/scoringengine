CERT_JWT ?= $(shell cat jwt.key)

scoringengine: *.go
	go build -v -o engine
	sudo setcap cap_net_bind_service,cap_net_raw+ep engine
	mv engine scoringengine
	pkill scoringengine || true

.PHONY: docker
docker:
	docker build -t scoringengine .
	docker create --name scoringengine scoringengine
	docker cp scoringengine:/scoringengine .
	docker rm -vf scoringengine
	upx -9 scoringengine

cert.crt cert.key:
	traefik-cert getcert -u certs.sprinkle.cloud -d scoreboard.netkoth.org -j "${CERT_JWT}" -c cert.crt -k cert.key

package: clean
	go generate
	go build -v -i -tags=embed

test: scoringengine
	./scoringengine

test-watch: scoringengine
	while true; do \
		./scoringengine || true; \
	done

.PHONY: clean
clean:
	rm -f scoringengine

watch:
	find *.go Makefile | entr make
