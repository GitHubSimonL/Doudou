#!/usr/bin/env python
#coding=utf8

"""
协议生成器。
"""
import sys, os, codecs


def gen_go_const(const_list):
    filepath = os.path.join(os.getcwd(), 'protocol', 'const.go')
    f = codecs.open(filepath, 'w+', 'utf-8')
    f.write("""package protocol\n\nconst (\n""")
    for item in const_list:
        if item[0] == '#':
            f.write("\n// %s\n" % (item[1:]))
            continue
        (code, name, desc) = item
        f.write('''%s = %s // %s\n''' % (name, code, desc))

    f.write(')\n')
    f.close()


# 10进制，16进制整数，浮点数
def is_number(s):
    try:
        int(s)
        return True
    except ValueError:
        try:
            int(s, 16)
            return True
        except ValueError:
            try:
                float(s)
                return True
            except ValueError:
                pass
    return False


def parse_const(buf):
    lines = [i.strip() for i in buf.split('\n')]

    d = []
    for line in lines:
        if not line: continue
        if line[0] == '#':
            d.append(line)
            continue
        el = line.split('-')
        if len(el) == 3:
            code, name, desc = el
        elif len(el) == 2:
            code, name = el
            desc = ""
        else:
            print('Error errcode:', el)

        if not is_number(code):
            code = '"' + code + '"'

        name = name.upper()
        d.append((code, name, desc))
    return d


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('usage: ./parse_const.py proto_dir [gen_dir]')
        sys.exit(0)

    path_pre = sys.argv[1]
    try:
        filepath = os.path.join(path_pre, 'const.txt')
        const_buf = codecs.open(filepath, 'r', 'utf-8').read()
    except IOError as e:
        print('Open proto file failed:', e)
        sys.exit(0)

    const_list = parse_const(const_buf)
    gen_go_const(const_list)
