#!/usr/bin/env python
# -*- coding: utf-8 -*-

from __future__ import print_function
import sys
import os
import csv
import codecs
import argparse


# 支持的字库分隔符
stop_punctuations = ['、', '，', ',']


def is_ignored_filename(filename):
    ignored_extension = [
    '~$',
    '-TNP-',
    ' - 副本',
    ]
    for text in ignored_extension:
        if filename.find(text) >= 0:
            return True
    return False


def find_stop_pos(line, start):
    idx = -1
    for ch in stop_punctuations:
        i = line.find(ch, start)
        if i >= 0:
            idx = i
            break
    return idx


def parse_line(table, line):
    start = 0
    while True:
        end = find_stop_pos(line, start)
        if end < 0:
            return
        word = line[start:end]
        start = end + 1
        table[word.strip()] = True


# 解析txt文件，单词以中文的逗号、顿号或者英文的逗号、换行符结尾
def parse_txt(table, filename):
    print('parse', filename)
    fp = codecs.open(filename, 'r', 'utf-8')
    lines = fp.readlines()
    fp.close()
    for line in lines:
        parse_line(table, line)


# 解析csv文件，格式与excel一致
def parse_csv(table, filename):
    f = codecs.open(filename, "r", "gbk") # excel另存为的csv默认是gbk
    reader = csv.reader(f, delimiter=",", quotechar="'")
    for row in reader:
        if len(row) == 2:
            word = row[1]
            table[word.strip()] = True


# 解析excel的.xls文件格式
def parse_xls(table, filename):
    import xlrd
    if is_ignored_filename(filename):
        return

    print('parse', filename)
    book = xlrd.open_workbook(filename, on_demand=True)
    sheet = book.get_sheet(0)
    for rx in range(sheet.nrows):
        row = sheet.row(rx)
        word = str(row[1].value)
        table[word.strip()] = True


# 解析excel的.xlsx文件格式
def parse_xlsx(table, filename):
    import openpyxl
    if is_ignored_filename(filename):
        return
    print('parse', filename)
    wb = openpyxl.load_workbook(filename, data_only=True, read_only=True)
    sheet = wb.active
    for i, sheet_row in enumerate(sheet.rows):
        cell = sheet_row[1]
        word = str(cell.value)
        table[word.strip()] = True
    wb.close()


def run(args):
    dictionary = {}  # 字库
    for root, dirs, files in os.walk(args.dir):
        for filename in files:
            if filename.endswith('.txt'):
                parse_txt(dictionary, root + '/' + filename)
            elif filename.endswith('.csv'):
                parse_csv(dictionary, root + '/' + filename)
            elif filename.endswith('.xls'):
                parse_xls(dictionary, root + '/' + filename)
            elif filename.endswith('.xlsx'):
                parse_xlsx(dictionary, root + '/' + filename)

    if len(dictionary) == 0:
        print('no bad words read')
        return

    print('%d bad words read to dictionary' % len(dictionary) )

    rows = [['A_STR_WORD']]  #  添加csv title
    for name in dictionary:
        rows.append([name])

    f = codecs.open(args.out, 'w', 'utf-8')
    w = csv.writer(f, delimiter=',', lineterminator='\n', quotechar='"', quoting=csv.QUOTE_MINIMAL)
    w.writerows(rows)
    f.close()
    print('wrote %d profanity words to %s' % (len(dictionary), args.out))


def main():
    parser = argparse.ArgumentParser(description="export profanity words")
    parser.add_argument("-d", "--dir", default="E:/Projects/server/data/屏蔽字库", help="directory of profanity files")
    parser.add_argument("-o", "--out", default="DirtyWords.csv", help="output file path")
    args = parser.parse_args()
    run(args)


if __name__ == '__main__':
    if sys.version.startswith('2'):
        reload(sys)
        sys.setdefaultencoding('utf-8')  # python2设置为utf-8
    main()

