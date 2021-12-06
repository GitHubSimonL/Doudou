#!/usr/bin/env python3
#coding: utf-8

import ply.lex as lex
import ply.yacc as yacc
import numpy as np
import pandas as pd
import functools
import sys
import os
import os.path

reserved = ('RANGE', 'UNIQUE', 'EMPTY', 'NOTEMPTY', 'AND', 'OR', 'IS', 'ANY',
            'NOT', 'STARTS', 'ENDS', 'REGEX', 'LENGTH', 'LOWER', 'UPPER',
            'REF', 'IF')
tokens = reserved + ('NAME', 'LPAREN', 'RPAREN', 'COLON', 'COMMA', 'PERIOD',
                     'SLASH', 'STR_LITERAL', 'NUM_LITERAL', 'WILDCARD',
                     'DOLLAR')
reserved_map = {}
for r in reserved:
    reserved_map[r.lower()] = r

t_LPAREN = r'\('
t_RPAREN = r'\)'
t_COLON = r':'
t_COMMA = r','
t_PERIOD = r'.'
t_SLASH = r'/'
t_STR_LITERAL = r'"[^"]*"'
t_WILDCARD = r'\*'
t_DOLLAR = r'\$'
t_ignore = " \t"


def t_NAME(t):
    r'[a-zA-Z_][a-zA-Z0-9_]*'
    t.type = reserved_map.get(t.value, "NAME")
    return t


def t_NUM_LITERAL(t):
    r'-?[0-9]+(\.[0-9]+)?'
    try:
        t.value = int(t.value)
    except ValueError:
        print("Integer value too large %d", t.value)
        t.value = 0
    return t


def t_newline(t):
    r'\n+'
    t.lexer.lineno += t.value.count("\n")


def t_error(t):
    print("Illegal character '%s'" % t.value[0])
    t.lexer.skip(1)


def series_or(sa, sb):
    return functools.reduce(lambda x, y: x | y, [sa, sb])


def series_and(sa, sb):
    return functools.reduce(lambda x, y: x & y, [sa, sb])


def series_op1(sa, sb):
    tmp = sa
    for i in range(len(sa)):
        if sa.iloc[i] and not sb.iloc[i]:
            tmp.iloc[i] = False
        else:
            tmp.iloc[i] = True
    return tmp


def series_op2(sa, sb, sc):
    tmp = sa
    for i in range(len(sa)):
        if sa.iloc[i] and sb.iloc[i]:
            tmp.iloc[i] = True
        elif not sa.iloc[i] and sc.iloc[i]:
            tmp.iloc[i] = True
        else:
            tmp.iloc[i] = False
    return tmp


def series_true(s):
    for idx, b in s.iteritems():
        if not b:
            return False
    return True


def check_result(ret):
    global validok, validrets

    bseries = ret[1]
    if not series_true(bseries):
        validok = False

    validrets.append((colname, ret))


def p_Scheam(p):
    '''Schema : Body'''
    check_result(p[1])


def p_Body(p):
    '''Body : BodyPart'''
    p[0] = p[1]


def p_BodyPart(p):
    'BodyPart : ColumnDefinition'
    p[0] = p[1]


def p_ColumnDefinition(p):
    'ColumnDefinition : ColumnIdentifier COLON ColumnRule'
    p[0] = p[3]


def p_ColumnIdentifier(p):
    'ColumnIdentifier : NAME'
    p[0] = p[1]


def p_ColumnRule(p):
    'ColumnRule : ColumnValidationExpr'
    p[0] = p[1]


def p_ColumnValidationExpr(p):
    '''ColumnValidationExpr : CombinatorialExpr
                            | NonCombinatorialExpr'''
    p[0] = p[1]


def p_CombinatorialExpr(p):
    '''CombinatorialExpr : OrExpr
                        | AndExpr'''
    p[0] = p[1]


def p_AndExpr(p):
    'AndExpr : NonCombinatorialExpr AND ColumnValidationExpr'
    left = (p[1])
    right = (p[3])
    ret = series_and(left[1], right[1])
    p[0] = (p[2], ret)


def p_OrExpr(p):
    'OrExpr : NonCombinatorialExpr OR ColumnValidationExpr'
    left = (p[1])
    right = (p[3])
    s = df[colname]
    ret = series_or(left[1], right[1])
    p[0] = (p[2], ret)


def p_NonCombinatorialExpr(p):
    '''NonCombinatorialExpr :  NonConditionalExpr
                            | ConditionalExpr'''

    p[0] = p[1]


def p_NonConditionalExpr(p):
    '''NonConditionalExpr : SingleExpr
                        | ExplicitSingleExpr
                        | ParenthesizedExpr'''
    p[0] = p[1]


def p_ConditionalExpr(p):
    '''ConditionalExpr : IfExpr'''

    p[0] = p[1]


def p_SingleExpr(p):
    '''SingleExpr : IsExpr
                    | AnyExpr
                    | NotExpr
                    | UniqueExpr
                    | EmptyExpr
                    | NotEmptyExpr
                    | StartsWithExpr
                    | EndsWithExpr
                    | RegExpExpr
                    | LengthExpr
                    | UpperCaseExpr
                    | LowerCaseExpr
                    | RefExpr
                    | RangeExpr'''
    p[0] = p[1]


def p_ExplicitSingleExpr(p):
    '''ExplicitSingleExpr : ExplicitContextExpr  SingleExpr'''

    global colname, tmpcol
    colname = tmpcol
    tmpcol = ''
    p[0] = p[2]


def p_ParenthesizedExpr(p):
    '''ParenthesizedExpr : LPAREN ColumnValidationExpr RPAREN'''
    p[0] = p[2]


def p_IsExpr(p):
    'IsExpr : IS LPAREN StringProvider RPAREN'

    ret = df[colname].isin([p[3]])
    p[0] = (p[1], ret)


def p_AnyExpr(p):
    '''AnyExpr :  ANY LPAREN StringProviderList RPAREN'''

    candidate = p[3]
    data = df[colname]
    ret = data.apply(lambda x: True if str(x) in candidate else False)
    p[0] = (p[1], ret)


def p_NotExpr(p):
    'NotExpr :   NOT LPAREN StringProvider RPAREN'

    ret = df[colname].isin([p[3]])
    ret = ret.apply(lambda x: not x)
    p[0] = (p[1], ret)


def p_RangeExpr(p):
    '''RangeExpr : RANGE LPAREN NumericOrAny COMMA NumericOrAny  RPAREN'''

    min, max = p[3], p[5]
    if min == '*':
        min = -sys.maxint - 1
    if max == '*':
        max = sys.maxint
    ret = df[colname].between(min, max)
    p[0] = (p[1], ret)


def p_UniqueExpr(p):
    '''UniqueExpr : UNIQUE'''

    data = df[colname]
    tmp = data.value_counts()
    f = lambda x: False if tmp[x] > 1 else True
    ret = data.apply(f)
    p[0] = (p[1], ret)


def p_EmptyExpr(p):
    '''EmptyExpr : EMPTY'''

    data = df[colname]
    ret = data.apply(lambda x: True if len(str(x)) == 0 else False)
    p[0] = (p[1], ret)


def p_NotEmptyExpr(p):
    '''NotEmptyExpr : NOTEMPTY'''

    data = df[colname]
    ret = data.apply(lambda x: False if len(str(x)) == 0 else True)
    p[0] = (p[1], ret)


def p_StartsWithExpr(p):
    '''StartsWithExpr : STARTS LPAREN StringProvider RPAREN'''

    data = df[colname]
    p[0] = (p[1], data.str.startswith(p[3]))


def p_EndsWithExpr(p):
    '''EndsWithExpr : ENDS LPAREN StringProvider RPAREN'''

    data = df[colname]
    p[0] = (p[1], data.str.endswith(p[3]))


def p_LowerCaseExpr(p):
    '''LowerCaseExpr : LOWER'''

    data = df[colname]
    ret = data.apply(lambda x: True if x.lower() == x else False)
    p[0] = (p[1], ret)


def p_UpperCaseExpr(p):
    '''UpperCaseExpr : UPPER'''

    data = df[colname]
    ret = data.apply(lambda x: True if x.upper() == x else False)
    p[0] = (p[1], ret)


def p_RegExpExpr(p):
    '''RegExpExpr : REGEX LPAREN StringProvider RPAREN'''

    data = df[colname]
    p[0] = (p[1], data.str.contains(p[3], regex=True))


def p_LengthExpr(p):
    '''LengthExpr : LENGTH LPAREN NumericLiteral COMMA NumericLiteral RPAREN'''

    data = df[colname].str.len()
    ret = data.apply(lambda x: True if x >= p[3] and x <= p[5] else False)
    p[0] = (p[1], ret)


def p_RefExpr(p):
    '''RefExpr : REF LPAREN DOLLAR NAME RPAREN
                | REF LPAREN DOLLAR NAME PERIOD NAME RPAREN'''
    if len(p) == 6:
        newcol = p[4]
        ret = df[colname].isin(df[newcol].tolist())
        p[0] = (p[1], ret)
    else:
        newtab = p[4]
        newcol = p[6]
        newdf = get_csv(os.path.join(basepath, newtab + ".csv"))
        ret = df[colname].isin(newdf[newcol].tolist())
        p[0] = (p[1], ret)


def p_ExplicitContextExpr(p):
    '''ExplicitContextExpr : ColumnRef SLASH'''

    p[0] = p[1]


def p_IfExpr(p):
    '''IfExpr : IF LPAREN ColumnValidationExpr COMMA SingleExpr COMMA SingleExpr RPAREN
                | IF LPAREN ColumnValidationExpr COMMA SingleExpr RPAREN'''

    if len(p) == 7:
        ret1 = p[3][1]
        ret2 = p[5][1]
        p[0] = (p[1], series_op1(ret1, ret2))
    else:
        ret1 = p[3][1]
        ret2 = p[5][1]
        ret3 = p[7][1]
        p[0] = (p[1], series_op2(ret1, ret2, ret3))


def p_ColumnRef(p):
    '''ColumnRef : DOLLAR ColumnIdentifier'''

    global colname, tmpcol
    tmpcol = colname
    colname = p[2]
    p[0] = p[2]


def p_NumericOrAny(p):
    '''NumericOrAny :	NumericLiteral
                    | WildcardLiteral'''
    p[0] = p[1]


def p_WildcardLiteral(p):
    '''WildcardLiteral : WILDCARD'''
    p[0] = p[1]


def p_StringProvider(p):
    'StringProvider : StringLiteral'
    p[0] = p[1].strip('\"')


def p_StringProviderList(p):
    '''StringProviderList : StringProvider
                        | StringProvider COMMA StringProviderList'''
    if len(p) == 2:
        p[0] = [p[1]]
    else:
        p[0] = [p[1]] + p[3]


def p_StringLiteral(p):
    'StringLiteral : STR_LITERAL'
    p[0] = p[1]


def p_NumericLiteral(p):
    'NumericLiteral : NUM_LITERAL'
    p[0] = p[1]


def p_error(p):
    print("Syntax error at '%s'" % p.value)


basepath = ''  #csv文件所在的路径
csvname = ''  #当前处理的csv文件名
colname = ''  #当前csv文件的列名
tmpcol = ''  #临时保存用
df = None  #当前处理的文件内容
dfcache = {}  #缓存已经读取进来的csv文件
maxcnt = 20  #单列错误最多显示条数
validok = True
validrets = []
lexer = lex.lex()
parser = yacc.yacc(debug=False)


def read_csv(csv):
    d = pd.read_csv(csv)
    d = d.replace(np.NaN, '')
    return d


def get_csv(csv):
    global dfcache

    if csv not in dfcache:
        dfcache[csv] = read_csv(csv)

    return dfcache[csv]


def validate_all(csvpath, schemapath):
    for fn in os.listdir(csvpath):
        if not os.path.isfile(
                os.path.join(csvpath, fn)) or not fn.endswith('.csv'):
            continue

        csv = os.path.join(csvpath, fn)
        schema = os.path.join(schemapath, fn + 's')
        if not os.path.isfile(csv) or not os.path.isfile(schema):
            continue
        validate_csv(csv, schema)


def validate_csv(csv, schema):
    global df, csvname, colname, basepath, validok, validrets

    validok = True
    validrets = []
    basepath = os.path.dirname(csv)
    csvname = os.path.basename(csv)
    df = get_csv(csv)

    tip = 'Validating {0} ...'.format(csvname)
    print(tip),
    for line in open(schema):
        line = line.split('#')[0].strip()
        if len(line) <= 0:
            continue
        colname = (line.split(':')[0]).strip()
        parser.parse(line)
    if validok:
        ret = '[PASS]'.rjust(120 - len(tip), ' ')
        print('\033[1;32m%s\033[0m' % ret)
    else:
        ret = '[FAIL]'.rjust(120 - len(tip), ' ')
        print('\033[1;31m%s\033[0m' % ret)

        for col, v in validrets:
            rule = v[0]
            cnt = 0
            for idx, b in v[1].iteritems():
                if not b:
                    cnt += 1
                    if cnt > maxcnt:
                        print(
                            '\t\033[1;33mcolumn %s too many errors, skipping...\033[0m'
                            % col)
                        break
                    row = idx + 2  #下标从0开始，并且不包含标题
                    print('\tcolumn: %s, row: %d, value: %s, rule: %s ?' %
                          (col, row, df[col].iloc[idx], rule))


if __name__ == '__main__':
    if len(sys.argv) < 3:
        print('Usage: %s <csv path> <schema path>' % (sys.argv[0]))
        exit(1)

    csvpath = sys.argv[1]
    if not os.path.isdir(csvpath):
        print('Error: bad csv path.')
        exit(2)

    schemapath = sys.argv[2]
    if not os.path.isdir(schemapath):
        print('Error: bad schema path.')
        exit(3)

    try:
        validate_all(csvpath, schemapath)
    except BaseException as e:
        print(e.message)
