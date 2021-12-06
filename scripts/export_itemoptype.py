#!/usr/bin/python3
# -*-coding:utf-8-*-

from __future__ import unicode_literals
from __future__ import print_function

import os
import codecs
import argparse


def parse_line(line):
    i = line.find("=")
    if i < 0:
        idx =line.find("//")

        if idx < 0:
            return "",False
        else:
            return """
# %s
-------
| ID | 名称  | 使用效果
| --- | ---- |-----|"""%line[idx+2:].strip(),True

    name = line[:i].strip()
    j = line.index(";", i+1)
    value = int(line[i+1:j].strip())
    if value == 0:
        return"",False
    k = line.find("//", j+1)
    comment = ''
    if k >= 0:
        comment = line[k+2:].strip()


    return "| %d | %s | %s" % (value, name, comment),False


def run(args):
    allStr = []
    latestIsTitle = False
    with codecs.open(args.input, "r", "utf-8") as f:
        in_itemop = False
        lines = f.readlines()

        for line in lines:
            if in_itemop and line.find("}") >= 0:
                in_itemop = False
                break
            if line.find("enum ItemOpType {") >= 0:
                in_itemop = True

            if in_itemop:
             tmpStr, isTitle = parse_line(line)
             if tmpStr == "":
                continue

             if isTitle and latestIsTitle:
                allStr = allStr[:len(allStr)-1]

             latestIsTitle = isTitle
             allStr.append(tmpStr)



    for tmpStr in allStr:
        print(tmpStr)


def main():
    parser = argparse.ArgumentParser(description="export item_op_type from proto")
    parser.add_argument("-i", "--input", default="../../protocol/enum.proto", help="input proto source file")
    args = parser.parse_args()
    print(os.getcwd())
    run(args)


if __name__ == '__main__':
    main()
