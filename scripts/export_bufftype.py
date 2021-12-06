#!/usr/bin/python3
# -*-coding:utf-8-*-

from __future__ import unicode_literals
from __future__ import print_function

import os
import csv
import codecs
import argparse


# 解析一行: `BuffTypeFoodProduct = 100; // 食物产量`
def parse_line(line):
    i = line.find("=")
    if i < 0:
        return
    name = line[:i].strip()
    j = line.index(";", i+1)
    value = int(line[i+1:j].strip())
    if value == 0:
        return
    k = line.find("//", j+1)
    comment = ''
    if k >= 0:
        comment = line[k+2:].strip()
    return {"name": name, "value": value, "comment": comment}


# 解析所有buff列表
def parse_buffs(filepath):
    buffs = []
    with codecs.open(filepath, "r", "utf-8") as f:
        in_buff = False
        lines = f.readlines()
        for line in lines:
            if in_buff:
                buff = parse_line(line)
                if buff is not None:
                    buffs.append(buff)
            if line.find("enum BuffType {") >= 0:
                in_buff = True
            if in_buff and line.find("}") >= 0:
                in_buff = False
                break
    return buffs


# 输出为markdown
def output_markdown(filepath, buffs):
    content = """
| ID | 名称  | Buff效果备注
| --- | ---- |-----|
"""
    for buff in buffs:
        content += "| %d | %s | %s\n" % (buff["value"], buff["name"], buff["comment"])
    f = codecs.open(filepath, 'w+', 'utf-8')
    f.write(content)
    f.close()
    print('saved to', filepath)


# 输出为csv
def output_csv(filepath, buffs):
    f = codecs.open(filepath, 'w+', 'utf-8')
    w = csv.writer(f, delimiter=',', lineterminator='\n', quotechar='"', quoting=csv.QUOTE_ALL)
    header = ['Buff ID', 'Buff名称', 'Buff效果备注']
    w.writerow(header)
    for buff in buffs:
        row = [buff["value"], buff["name"], buff["comment"]]
        w.writerow(row)
    f.close()
    print('saved to', filepath)


def run(args):
    buffs = parse_buffs(args.input)
    if args.output.endswith('.md'):
        output_markdown(args.output, buffs)
    elif args.output.endswith('.csv'):
        output_csv(args.output, buffs)
    else:
        raise RuntimeError("unsupported file extension")


def main():
    parser = argparse.ArgumentParser(description="export buff types from proto")
    parser.add_argument("-i", "--input", default="../../protocol/enum.proto", help="input proto source file")
    parser.add_argument("-o", "--output", default="buff_type.csv", help="output file path")
    args = parser.parse_args()
    print(os.getcwd())
    run(args)


if __name__ == '__main__':
    main()
