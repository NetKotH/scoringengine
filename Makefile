CERT_JWT ?= eyJhbGciOiJSUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJjZXJ0Ijp7ImRvbWFpbnMiOlsic2NvcmVib2FyZC5uZXRrb3RoLm9yZyJdfSwiZXhwIjoxNTc4MDA1NDQ5LCJpYXQiOjE1NzAyMjk0NDksIm5iZiI6MTU3MDIyOTQ0OX0.ZdUYaCxVAWdV1UuQ9S-POFEODnCiat-HPr2OuEuAUOUa6sD8axg_zRF8SY-V-SDAYdFfTWWBLEGblsQP-iHGE0WqTmxD8AMGLNc_RpaRdAcSRZXCxLaQ6yDxu7YcEe5FiSZM5epY4u2qqcvQEy8LKpIA1OQvnqmG6z-HaYnbZ1g0oIZSpvOwbQtM_6rCcDgsxJxFP29w81Um32jUkAH65rkqMl6aQbN0aUw8gp-_L9L-Sq70JEykSZR_hNkDw6rQcEoBictuaKyEIN3IY9cAcZS_q0q6FxVLrbIQ0GG1iTBvzTq1CQd-ykCEJtCSARk9CY-Tn07Hu-xMzSJX1JtMtQ

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
