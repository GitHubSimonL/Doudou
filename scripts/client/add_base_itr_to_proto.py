#!/usr/bin/env python
# coding=utf-8
import getopt
import re
import sys

inPaths = []


def usage():
    print("add_base_itr_to_proto.py: \n"
          "args: \n"
          " --in path   eg: protocol/demo1.pb.go\n"
          "eg: \n"
          " add_base_itr_to_proto.py --in protocol/demo1.pb.go")


def get_opt():
    try:
        options, args = getopt.getopt(sys.argv[1:], "", ["in=", "help"])
        for name, value in options:
            if name == '--in':
                global inPaths
                inPaths.append(value)
            if name == '--help':
                usage()
                return False

        return True
    except Exception as e:
        print("get_opt error: %s" % (e))
        usage()
        return False


def sedFile():
    print("sed ====")
    for fname in inPaths:
        fp = open(fname, "r")
        lines = []

        for line in fp:
            lines.append(line)

        fp.close()

        for index, line in enumerate(lines):
            ret0 = re.match("^\s*type\s+[a-zA-Z_][a-zA-Z0-9_]*Req\s*struct\s*{|"
                            "^\s*type\s+[a-zA-Z_][a-zA-Z0-9_]*Ack\s*struct\s*{|"
                            "^\s*type\s+[a-zA-Z_][a-zA-Z0-9_]*Ntf\s*struct\s*{|", line)
            if ret0.group():
               print(line)
#                 ret1 = re.findall("[a-zA-Z_][a-zA-Z0-9_]*", ret0.group())




#                 if ret1:
#                     s = "Processor.Register((*" + ret1[1] + ")(nil))\n"
#                     outFp.writelines(s)
#                        print(ret1[1])

#     outFp.writelines("}\n")


def gen():
    if get_opt() == False:
        return
    sedFile()


if __name__ == '__main__':
    gen()
