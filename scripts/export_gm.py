#!/usr/bin/python3
# -*-coding:utf-8-*-

from __future__ import unicode_literals
from __future__ import print_function

import os
import sys
import csv
import codecs
import argparse
import json
import time

import antlr4

SCRIPT_PATH = os.path.join(os.path.dirname(__file__))
sys.path.append('./')

from parsers.GoLexer import GoLexer
from parsers.GoParser import GoParser
from parsers.GoParserListener import GoParserListener


pkg_name_mapping = {
    'gs': '游戏服',
    'gate': '网关服',
    'login': '登录服',
    'center': '联盟中心服',
    'mgard': '内城服',
    'maild': '邮件服',
    'gma': 'GM服',
    'global_activity': '全服活动',
    'multi_map': '多人副本',
}


class MyGMFuncListener(GoParserListener):
    def __init__(self, content):
        self.content = content
        self.fnname = ''
        self.args = []

    def enterMethodDecl(self, ctx: GoParser.FunctionDeclContext):
        self.fnname = str(ctx.IDENTIFIER())
        params = ctx.signature().parameters()
        for param in params.parameterDecl():
            type_ctx = param.type_()
            ellipsis = param.ELLIPSIS() or ''
            if type_ctx.typeName() is None:
                type_ctx = param.type_().typeLit()
                i = type_ctx.start.start
                j = type_ctx.stop.stop
                type_name = str(ellipsis) + self.content[i:j+1]
            else:
                type_name = str(ellipsis) + str(type_ctx.typeName().IDENTIFIER())

            ident_list = param.identifierList().IDENTIFIER()
            for ident in ident_list:
                argument = (type_name, str(ident))
                self.args.append(argument)


# 对于一些函数签名格式（如下划线表示的空白变量，可变个数参数，以及多个相同类型的变量类型可以省略等等）
# 这里为了解析的准确性，使用antlr来解析函数签名
def parse_func_signature(content) -> tuple:
    input_stream = antlr4.InputStream(content)
    lexer = GoLexer(input_stream)
    parser = GoParser(antlr4.CommonTokenStream(lexer))
    tree = parser.methodDecl()
    walker = antlr4.ParseTreeWalker()
    listener = MyGMFuncListener(content)
    walker.walk(listener, tree)
    return (listener.fnname, listener.args)


def parse_file(filepath, class_name, gm_funs):
    print('start parse', filepath)
    f = codecs.open(filepath, 'r', 'utf-8')
    lines = f.readlines()
    f.close()

    idx = 0
    pkg_name = ''
    while idx < len(lines):
        start_idx = idx
        line = lines[idx].strip()
        idx += 1
        if line.startswith("package "):
            pkg_name = line[8:].strip()
        if not line.startswith("func "):
            continue
        while line.endswith(","):  # 处理跨越多行的参数
            line += lines[idx].strip()
            idx += 1
        start = line.find(class_name + ')')
        if start > 0:
            start = start + len(class_name + ')')
            # print(line)
            fnsig = line[start: line.index(')', start) + 1].strip()
            i = fnsig.index('(')
            fnname = fnsig[:i]
            if not fnname[0].isupper():  # 首字母大写才是可导出的Go函数
                continue
            last_line = lines[start_idx - 1].strip()
            comment = ''
            if last_line.startswith('//'):
                comment = last_line[2:].strip()
                if comment.startswith(fnname):
                    comment = comment[len(fnname):].strip()
            (fnname, args) = parse_func_signature(line + '\n}')
            item = {
                'name': fnname,
                'comment': comment,
                'args': args,
            }
            if pkg_name not in gm_funs:
                gm_funs[pkg_name] = []
            gm_funs[pkg_name].append(item)

    return gm_funs


def format_func_args(args) -> str :
    s = ''
    i = 0
    for tup in args:
        if tup[1] == '_':
            s += tup[0]
        else:
            s += '%s %s' % (tup[0], tup[1])
        if i + 1 < len(args):
            s += ', '
        i += 1
    return s


def format_func_style(name, args) -> str:
    s = name + ('(')
    i = 0
    for arg in args:
        s += arg[1]
        if i + 1 < len(args):
            s += ', '
        i += 1
    return s + ')'


def output_markdown(filepath, gm_funs):
    content = """
| 模块 | GM格式  | 参数类型 | 备注
| --- | ---- |-----|----
"""
    for pkg_name in gm_funs:
        for v in gm_funs[pkg_name]:
            func_style = format_func_style(v['name'], v['args'])
            args_style = format_func_args(v['args'])
            content += "| %s | %s | %s | %s\n" % (pkg_name_mapping[pkg_name], func_style, args_style, v["comment"])

    f = codecs.open(filepath, 'w+', 'utf-8')
    f.write(content)
    f.close()
    print('saved to', filepath)


def output_csv(filepath, gm_funs):
    f = codecs.open(filepath, 'w+', 'utf-8')
    w = csv.writer(f, delimiter=',', lineterminator='\n', quotechar='"', quoting=csv.QUOTE_ALL)
    header = ['模块', 'GM格式', '备注', '参数类型']
    w.writerow(header)
    for pkg_name in gm_funs:
        for v in gm_funs[pkg_name]:
            func_style = format_func_style(v['name'], v['args'])
            args_style = format_func_args(v['args'])
            row = [pkg_name_mapping[pkg_name], func_style, v['comment'], args_style]
            w.writerow(row)
    f.close()
    print('saved to', filepath)


def run(args):
    base_gm_funs = {}
    mod_gm_funs = {}
    dir = os.path.abspath(args.dir)
    print('start looking dir:', dir)
    common_handle_file = os.path.join('internal', 'gm_module', 'handler.go')
    for root, dirs, filenames in os.walk(dir):
        for filename in filenames:
            if root.find('vendor') > 0 or root.find('scripts') > 0:
                continue
            filepath = os.path.join(root, filename)
            filepath = os.path.abspath(filepath)
            if filepath.endswith(common_handle_file):
                parse_file(filepath, 'CommonGMHandler', base_gm_funs)
            if filepath.find('gm_func') > 0 and filepath.endswith('.go'):
                parse_file(filepath, 'gmHandler', mod_gm_funs)


    if len(mod_gm_funs) == 0:
        print('no GM command parsed')
        return

    # 把公共GM命令合并到每一个module内
    for pkg_name in mod_gm_funs:
        mod_gm_funs[pkg_name] += base_gm_funs['gm_module']

    if args.output == 'console':
        print(json.dumps(mod_gm_funs, ensure_ascii=False))
    elif args.output.endswith('.md'):
        output_markdown(args.output, mod_gm_funs)
    elif args.output.endswith('.csv'):
        output_csv(args.output, mod_gm_funs)
    else:
        raise RuntimeError("unsupported file extension")


def main():
    parser = argparse.ArgumentParser(description="export GM from go source")
    parser.add_argument("-d", "--dir", default="../", help="Go module directory")
    parser.add_argument("-o", "--output", default="gm.csv", help="output file path")
    args = parser.parse_args()
    run(args)


if __name__ == '__main__':
    main()
