#!/usr/bin/env python
# coding=utf-8
import getopt
import re
import sys

inpaths = []
outfile = ""


def usage():
    print("gen_req_proto.py: gen msg.go\n"
          "args: \n"
          " --in path   eg: protocol/api.proto\n"
          " --out file  eg: src/protobuf/msg.go\n"
          "eg: \n"
          " gen_req_proto.py --in protocol/api.proto --in protocol/internal_api.proto --out protocol/msg.go")


def get_opt():
    try:
        options, args = getopt.getopt(sys.argv[1:], "", ["in=", "out=", "help"])
        for name, value in options:
            if name == '--in':
                global inpaths
                inpaths.append(value)
            if name == '--out':
                global outfile
                outfile = value
            if name == '--help':
                usage()
                return False

        if outfile == "":
            return False

        return True
    except Exception as e:
        print("get_opt error: %s" % (e))
        usage()
        return False


def gen_file():
    outFp = open(outfile, "w")

    outFp.writelines("// generate by gen_msg.gp.py . DO NOT EDIT\n"
                     "\n"
                     "package protocol\n"
                     "\n"
                     "import \"server/base/network/parser\"\n"
                     "\n"
                     "// Processor proto\n"
                     "var Processor = parser.NewProcessor()\n"
                     "// Init register msg\n"
                     "func init () {\n"
                     )

    for fname in inpaths:
        fp = open(fname, "r")
        outFp.writelines("// " + fname + "\n")
        for line in fp:
            # message xxxReq {
            ret0 = re.match("^\s*message\s+[a-zA-Z_][a-zA-Z0-9_]*Req\s*{|"
                            "^\s*message\s+[a-zA-Z_][a-zA-Z0-9_]*Ack\s*{|"
                            "^\s*message\s+[a-zA-Z_][a-zA-Z0-9_]*Ntf\s*{|", line)

            if ret0:
                ret1 = re.findall("[a-zA-Z_][a-zA-Z0-9_]*", ret0.group())

                if ret1:
                    s = "Processor.Register((*" + ret1[1] + ")(nil))\n"
                    outFp.writelines(s)

    outFp.writelines("}\n")


def gen():
    if get_opt() == False:
        return
    gen_file()


if __name__ == '__main__':
    gen()
