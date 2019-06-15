.PHONY:  default  test  test-docker

default: test

test:
	@scripts/run-all-tests
	@echo ========================================
	@git grep TODO  -- '**.elm' || true
	@git grep FIXME -- '**.elm' || true

test-docker:
	docker-compose --version
	docker-compose build --pull elm
	docker-compose up --abort-on-container-exit --exit-code-from=elm --force-recreate
