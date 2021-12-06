#!/usr/bin/env python
# coding=utf8

from __future__ import print_function
from __future__ import with_statement
import math
import xlrd
import sys
import os
import csv
import codecs
import traceback

LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'

class SimpleProgressBar():
    def __init__(self, header, width=50):
        self.last_x = -1
        self.width = width
        self.header = header

    def update(self, x):
        assert 0 <= x <= 100 # `x`: progress in percent ( between 0 and 100)
        if self.last_x == int(x): return
        self.last_x = int(x)
        pointer = int(self.width * (x / 100.0))
        sys.stdout.write( '\r%s : %d%% [%s]' % (self.header, int(x), '#' * pointer + '.' * (self.width - pointer)))
        sys.stdout.flush()
        if x == 100: print('')


def excel_style(row, col):
    """ Convert given row and column number to an Excel-style cell name. """
    result = []
    while col:
        col, rem = divmod(col - 1, 26)
        result[:0] = LETTERS[rem]
    return ''.join(result) + str(row)


# excel表格配置格式
#   `A_`开头，服务器、客户端都需要读取的列
#   `C_`开头，仅客户端需要读取的列
#   `S_`开头，仅服务器需要读取的列
#   `#`开头，注释列
def parse(in_file, out_file):
    try:
        data = xlrd.open_workbook(in_file)
        table = data.sheets()[0]
        nrows = table.nrows
        ncols = table.ncols
        rows = []
        row = []
        headers = {}
        for col in range(ncols):
            v = str(table.cell(0, col).value).strip()
            if v.startswith("C_") or v.startswith("#"):  # ignore this column
                continue
            if v.startswith("A_") or v.startswith("S_"):
                v = v[2:]
            strs = v.split('_')
            if len(strs) > 1:
                headers[col] = strs[0].upper().encode('utf-8')
            else:
                headers[col] = ""
            row.append(v)
        rows.append(row)
        #print headers
        #pb = SimpleProgressBar('rows')
        for i in range(nrows):
            #pb.update((i+1)*100/(nrows))
            if i == 0:
                continue
            row = []
            for j in range(ncols):
                v = table.cell(i, j).value
                t = headers.get(j)
                if t is None:
                    continue
                if t == "INT":
                    ct = table.cell(i, j).ctype     # cell type
                    if ct == xlrd.XL_CELL_NUMBER:
                        v = str(math.trunc(v))
                    elif ct == xlrd.XL_CELL_TEXT:
                        v = int(v.strip())  # test if cell contains integer value
                    elif ct == xlrd.XL_CELL_EMPTY:
                        v = v.strip()
                    else:
                        print(
                            '\033[1;35mfile %s format error, expect number at cell(%d, %d) %s \033[0m'
                            % (in_file, i + 1, j, excel_style(i + 1, j + 1)))
                        sys.exit(1)
                elif t == "STR":
                    if table.cell(i, j).ctype == xlrd.XL_CELL_NUMBER:
                        v = str(v.strip())
                else:
                    if table.cell(i, j).ctype == xlrd.XL_CELL_NUMBER:
                        v = str(math.trunc(v))
                v = str(v)
                row.append(v.encode("utf-8"))
            rows.append(row)
        with codecs.open(out_file, "w", "utf-8") as fp:
            writer = csv.writer(fp, delimiter=',', lineterminator='\n', quotechar='"', quoting=csv.QUOTE_ALL)
            writer.writerows(rows)
    except Exception as ex:
        print(ex)
        traceback.print_exc()


def scan(src_path, dst_path):
    if not os.path.isdir(src_path):
        print("Path error: ", src_path)
        return

    fullpath = os.path.split(os.path.realpath(__file__))[0]
    inlist = getlinelist(fullpath + "/convert-in.txt")
    exlist = getlinelist(fullpath + "/convert-ex.txt")

    for f in os.listdir(src_path):
        src_filename = os.path.join(src_path, f)
        if not os.path.isfile(src_filename):
            continue
        if src_filename.find("~") >= 0:    # ignore temp file
            continue
        name, ext = os.path.splitext(os.path.basename(f))
        fn = name + ext

        if (fn in inlist) or (len(inlist) == 0):
            if fn in exlist:
                print("WARN: ignore file %s" % (fn))
                continue
            if name != "battleaffect" and ext == ".xls":
                name += ".csv"
                dst_filename = os.path.join(dst_path, name)
                print("%s ==> %s" % (src_filename, dst_filename))
                parse(src_filename, dst_filename)


def getlinelist(fn):
    lines = []
    if os.path.isfile(fn):
        f = open(fn)
        while 1:
            line = f.readline()
            if not line:
                break
            lines.append(line.strip('\n'))
        f.close()
    return lines


if __name__ == "__main__":
    reload(sys)
    sys.setdefaultencoding('utf8')
    if len(sys.argv) != 3:
        print("Usage:")
        print("\t%s src_path dst_path" % (sys.argv[0]))
        sys.exit(0)
    try:
        scan(sys.argv[1], sys.argv[2])
    except Exception as ex:
        print(ex)
        traceback.print_exc()
