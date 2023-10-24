VERSION=`cat ./VERSION`

.PHONY: docker.test docker.prune hooks.style hooks.version-update hooks.setup update.dep

# build images for deploy
docker.test:
	docker build -t opabinia-test:${VERSION} -f build/test.Dockerfile .
	docker rmi opabinia-test:${VERSION}

docker.prune:
	docker image prune -f

# checks before commit
hooks.style:
	gofumpt -l -w . && go vet $$(go list ./... | grep -v /vendor/)

hooks.version-update:
	scripts/version_update.sh

hooks.setup:
	@cp scripts/pre-commit-hooks.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	@cp scripts/post-commit-hooks.sh .git/hooks/post-commit
	chmod +x .git/hooks/post-commit

update.dep:
	@go get -u ./...