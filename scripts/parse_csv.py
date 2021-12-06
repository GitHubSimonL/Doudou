#!/usr/bin/env python
# -*- coding: utf-8 -*-

from pyExcelerator import *
import sys
import string
import os

if len(sys.argv) == 3 :
    gd_dir = sys.argv[1]
    project_dir = sys.argv[2]

export_list = [
    {"path": "building.xlt", "export":{"BuildingCost": "building.csv"}},
    ]


def export_csv(xls_path, output_dir, xls_table):
    for sheet_name, values in parse_xls(xls_path, "cp1251"):
        print sheet_name, values
        key = sheet_name.encode("utf8", "backslashreplace")
        print sheet_name, values, key
        if not key in xls_table.keys():
            continue
        target = xls_table[key]
        print "Try to export %s to %s." % (key, target)
        matrix = [[]]
        for row_idx, col_idx in sorted(values.keys()):
            v = values[(row_idx, col_idx)]

            if isinstance(v, unicode):
                v = v.encode("utf8", "backslashreplace")
            else:
                v = str(v)
            v = v.replace(" ", "").replace("\r", "").replace("\n", "")
            try:
                v = (int)(string.atof(v))
            except:
                v = str(v)
            v = str(v)
            if v == "": puts("should not be null")
            last_row, last_col = len(matrix), len(matrix[-1])
            while last_row-1 < row_idx:
                matrix.extend([[]])
                last_row = len(matrix)
            while last_col < col_idx:
                matrix[-1].extend([''])
                last_col = len(matrix[-1])

            matrix[-1].extend([v])

        f = open(os.path.join(output_dir, target), "w")
        for row in matrix:
            csv_row = ",".join(row)
            f.write(csv_row + "\n")
        f.close()


def export_all(xls_dir, output_dir):
    for export_item in export_list:
        xls_path = os.path.join(xls_dir, export_item["path"])
        export_csv(xls_path, output_dir, export_item["export"])


def usage():
    print 'usage: ./parse_csv.py xls_dir output_dir'


if __name__ == '__main__':
    if len(sys.argv) < 3:
        usage()
        sys.exit()
    export_all(sys.argv[1], sys.argv[2])

