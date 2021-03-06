.PHONY: force test

SRC := $(wildcard *.ml) lexer.mll parser.mly opcode.ml cprim.ml
C := $(addprefix runtime/,main.c compare.c error.c instruct.c io.c main.c prim.c str.c)
HD := $(addprefix runtime/,common.h error.h instruct.h io.h jumptable.h prim.h str.h value.h)
CYAN := '\e[1;36m'
GREEN := '\e[1;32m'
RST := '\e[0m'

all: cc od run

cc: main.byte
	ln -sf $< $@

od: objdump.byte
	ln -sf $< $@

%.byte: $(SRC)
	ocamlbuild $@
	touch $@

run: $(C) $(HD)
	$(LINK.c) -I runtime $(filter %.c,$^) -o $@

runtime/jumptable.h: runtime/instruct.h
	sed -rn 's/([[:upper:]]+)/\&\&lbl_\1/;T;p' $< > $@

runtime/instruct.c: runtime/instruct.h
	{ echo 'const char *name_of_instructions[] = {'; sed -rn 's/([[:upper:][:digit:]]+).*/"\1",/;T;p' $<; echo '};';} > $@

opcode.ml: runtime/instruct.h
	awk '/[[:upper:]]/{sub(",","");print "let op"$$1"="n++;s[n-1]=$$1}END{printf "let name_of_opcodes=[|";for(i=0;i<n;i++)printf "\""s[i]"\";";print"|]"}' $< > $@

cprim.ml: runtime/prim.c
	awk '/cprims/{a=1} /\w+,/&&a==1{sub(",","");s[n++]=$$1} END{printf "let name_of_prims=[|";for(i=0;i<n;i++)printf "\""s[i]"\";";print"|]"}' $< > $@

test-%: force
	@echo -e $(CYAN)test $*$(RST)
	@t=$$(mktemp); ./cc tests/$*.ml -o $$t && ./run $$t && echo; r=$$?; $(RM) tests/$*.zo $$t; exit $$r

#test: $(addprefix test-,$(patsubst tests/%.ml,%,$(wildcard tests/*.ml)))
test:
	@all=0; fail=0; \
	for t in $(patsubst tests/%.ml,%,$(wildcard tests/*.ml)); do \
		$(MAKE) --no-print-directory test-$$t || ((fail++)); \
		((all++)); \
	done; \
	echo -e $(GREEN)$$fail/$$all cases failed$(RST)
