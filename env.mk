REVISION=$(shell ./scripts/gen_version.sh)

PB_SRC_DIR = ./proto
PB_GO_SRC_DIR = ./protocol
ALL_PB_SRC = $(wildcard $(PB_SRC_DIR:%=%/*.proto))
ALL_PB_GO_SRC = $(wildcard $(PB_GO_SRC_DIR:%=%/*.pb.go))
ALL_PB_SRC_OPT = $(subst ${PB_SRC_DIR}, --in ${PB_SRC_DIR}, ${ALL_PB_SRC})

