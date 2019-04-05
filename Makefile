CERT_JWT ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJjZXJ0Ijp7ImRvbWFpbnMiOlsic2NvcmVib2FyZC5uZXRrb3RoLm9yZyJdfSwiZXhwIjoxNTYyMTU5MjA4LCJpYXQiOjE1NTQzODMyMDgsIm5iZiI6MTU1NDM4MzIwOH0.wdA-F4Fq6qa6TgztWifXO5tZEHl4R1EijF7scHSF7m4oFuIy8UhOprGvzZcPYaEF1MeMdGk34F0jjoFJiIoEhoRgEedRVAA0OswScpHDj8dLlro80w751AFEQeeRB8jXj64eY_uQryfS6rt7kQFV2VlzOOkEkll6vXFGq3hpwvXMeusKMRsC-pfhHQt7-nLu6-PR84KWEmXCO6HmIwb-mgTfOYAs21GODolIIDYzXnGZCHqJmgDRNn-j3qBAkYX-vLqaOWq6EVc1wiB2rcqe1ExkaZ3Z3qhX7Yol9x0xaOZLzyCQ3DsmWotOc2eT7XzJJu797klxjHmllU997BvVBw

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
