NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

OK_STRING=$(OK_COLOR)[OK]$(NO_COLOR)
ERROR_STRING=$(ERROR_COLOR)[ERRORS]$(NO_COLOR)
WARN_STRING=$(WARN_COLOR)[WARNINGS]$(NO_COLOR)

AWK_CMD = awk '{ printf "%-30s %-10s\n",$$1, $$2; }'
PRINT_ERROR = printf "$@ $(ERROR_STRING)\n" | $(AWK_CMD) && printf "$(CMD)\n$$LOG\n" && false
PRINT_WARNING = printf "$@ $(WARN_STRING)\n" | $(AWK_CMD) && printf "$(CMD)\n$$LOG\n"
PRINT_OK = printf "$@ $(OK_STRING)\n" | $(AWK_CMD)
BUILD_CMD=\
	LOG=$$($(CMD) 2>&1);\
	if [ $$? -eq 1 ]; then $(PRINT_ERROR);\
	elif [ "$$LOG" != "" ] ; then $(PRINT_WARNING);\
	else $(PRINT_OK); fi;

all: build ut e2e

build:
	@@echo "$(OK_COLOR)"'building>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@# In order to run this successfully, add this line:
	@# build --genrule_strategy=standalone --spawn_strategy=standalone
	@# to  ~/.bazelrc;
	@# or wait for ease-gateway to fix.
	@bazel build --embed_label="$${BUILD_EMBED_LABEL}" ...
	@bazel run cmd/stress:latest
	@bazel run cmd/demo:latest
	@bazel run demo:latest
	@$(BUILD_CMD)

ut:
	@@echo "$(OK_COLOR)"'running unit tests>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@bazel test --test_size_filters="small" ...
	@$(BUILD_CMD)

serve:
	@@echo "$(OK_COLOR)"'starting services>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@bazel run cmd/stress:latest
	@docker-compose -f docker-compose-int-test.yml up -d --remove-orphans
	@@echo "$(OK_COLOR)"'waiting services...>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@while ! nc -q 1 localhost 6666 </dev/null; do sleep 1; done
	@$(BUILD_CMD)

e2e: serve
	@@echo "$(OK_COLOR)"'running e2e tests>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@bazel test --test_size_filters="medium" ...
	@$(BUILD_CMD)

stop:
	@@echo "$(OK_COLOR)"'stopping services>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@docker-compose -f docker-compose-int-test.yml down
	@$(BUILD_CMD)

verify:
	@@echo "$(OK_COLOR)"'Verify basic features >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>' && echo "$(NO_COLOR)"
	@docker-compose -f docker-compose-features.yml -f ../docker-compose/dev/skylb/docker-compose.yml up
	# Now manually check console and make sure there are no errors like panic, exceptions, etc. and press ctrl-c to stop.
	# TODO(fuyc): detect failures automatically.

.PHONY: build ut e2e
