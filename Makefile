CERT_JWT ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJjZXJ0Ijp7ImRvbWFpbnMiOlsic2NvcmVib2FyZC5uZXRrb3RoLm9yZyJdfSwiZXhwIjoxNTMzMDk2MDAwLCJpYXQiOjE1MzA1NzI2MjcsIm5iZiI6MTUzMDU3MjYyN30.zVTvFzlf9vG2OigCSYU2e6a-b6yz9dJrD6LFV-vAELFhvB3U7F_UpNLTyQY1EKaWzDHJJG4xLu0-cDYERoUnOePLXkHcAB8RAq3Lj-l3nEp8eYlyUAeI5v5d3HhS6Y8HvgnbuOVSo-8EIhZSUHwRpRSmH1dVr5lBnLov8D7_uExcXf1xVSH9WD8jgz0Twz_usaBeNdrXbghVHEdCYuCbS-5MbfUPUR61tkBgk3Y2WDnPY7SIiQDTWdqw9YVpMUjJrXJF0kFZadvrHljczGOOO-umHzkmxkppjSEGyFPxtjQOvU8pbN7nF5040vohGu-0yUzxsSfzEfOiJgIPaplb3w

scoringengine: *.go
	go build -v -i -o engine
	sudo setcap cap_net_bind_service,cap_net_raw+ep engine
	mv engine scoringengine
	#pkill -f scoringengine || true

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
	rm scoringengine

watch:
	find *.go Makefile | entr make
