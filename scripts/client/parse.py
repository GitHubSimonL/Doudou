#!/usr/bin/env python
#coding=utf8

"""
协议生成器。
"""
import sys, os

go_gen = False

class Api(object):
    def __init__(self, t, n, d, p, f):
        self.type = t
        self.name = n
        self.desc = d
        self.payload = p
        self.flag = f
        self.check()
    def __repr__(self):
        return "\nAPI type: %s \nname: %s desc: %s payload: %s\n" %\
            (self.type, self.name, self.desc, self.payload)

    def check(self):
        for i in [self.name, self.type, self.payload]:
            if (i is None) or (i == ''):
                print("API error:", self.type, self.name, self.payload)
                sys.exit(-1)
        if self.desc is None:
            self.desc = ''


class ProtoField(object):
    def __init__(self, n, t, array=False):
        self.name = n
        if t == 'int32':
            t = 's32'
        elif t == 'int64':
            t = 's64'
        elif t == 'int16':
            t = 's16'
        elif t == 'uint32':
            t = 'u32'
        elif t == 'uint64':
            t = 'u64'
        elif t == 'uint16':
            t = 'u16'
        elif t == 'bool':
            t = 'bool'
        elif t == 'int8':
            t = 's8'
        elif t == 'uint8':
            t = 'u8'
        self.type = t
        self.array = array
    def __repr__(self):
        return "[%s] type: %s, array: %s\n" % (self.name, self.type, self.array)

    def go_type(self):
        if self.type == 's32':
            return 'int32'
        elif self.type == 's64':
            return 'int64'
        elif self.type == 's16':
            return 'int16'
        elif self.type == 's8':
            return 'int8'
        elif self.type == 'u32':
            return 'uint32'
        elif self.type == 'u64':
            return 'uint64'
        elif self.type == 'u16':
            return 'uint16'
        elif self.type == 'u8':
            return 'uint8'

        if self.is_basic_type():
            return self.type
        else:
            return 'PKT_' + self.type

    def go_type_func(self):
        if self.type == 's32':
            return 'ReadS32'
        elif self.type == 's64':
            return 'ReadS64'
        elif self.type == 's16':
            return 'ReadS16'
        elif self.type == 'u32':
            return 'ReadU32'
        elif self.type == 's16':
            return 'ReadU16'
        elif self.type == 'string':
            return 'ReadString'
        elif self.type == 'bool':
            return 'ReadBool'
        elif self.type == 'rawdata':
            return 'ReadRawData'
        elif self.type == 's8':
            return 'ReadS8'
        elif self.type == 'u8':
            return 'ReadU8'
        return self.type

    def is_basic_type(self):
        return self.type in ['s32', 'u32', 's16', 'u16', 'string', 'bool', 's64', 's8', 'u8']


class Proto(object):
    def __init__(self, n, l):
        self.name = n
        self.fields = []
        for p in l:
            if len(p) == 2:
                self.fields.append(ProtoField(p[0], p[1]))
            elif len(p) == 3 and p[1] == 'array':
                self.fields.append(ProtoField(p[0], p[2], True))
            else:
                print("Bad proto type:", p)
                sys.exit(-1)

    def __repr__(self):
        return "Proto: %s\n%s" %\
            (self.name, ''.join([i.__repr__() for i in self.fields]))

    def check(self, proto_list):
        pass

def new_api(api_dict, api_type, api_name, api_desc, api_payload,api_flag):
    if api_type in api_dict:
        print("API dup: ", api_type)
        sys.exit(-1)
    api = Api(api_type, api_name, api_desc, api_payload, api_flag)
    api_dict[api_type] = api

def parse_api(api_buf,api_flag):
    L = [line.strip() for line in api_buf.split('\n')]
    L = [line for line in L if line and line[0] != '#']

    api_dict = {}

    start_api = False
    api_type = None
    api_name = None
    api_desc = None
    api_payload = None
    for line in L:
        idx = line.find(':')
        if idx < 0: continue
        if line[:idx] == 'packet_type':
            # 到新的API了，那么分析之前的API
            if start_api:
                new_api(api_dict, api_type, api_name, api_desc, api_payload,api_flag)
                api_type = api_name = api_desc = api_payload = None
            start_api = True
            api_type = int(line[idx+1:])
        elif line[:idx] == 'name':
            api_name = line[idx+1:]
        elif line[:idx] == 'payload':
            api_payload = line[idx+1:]
        elif line[:idx] == 'packet_type':
            api_type = True
        elif line[:idx] == 'desc':
            api_desc = line[idx+1:]
    if api_type:
        new_api(api_dict, api_type, api_name, api_desc, api_payload,api_flag)
    return api_dict

proto_list = []
def new_proto(proto_dict, proto_name, p_list):
    if proto_name in proto_dict:
        print("Proto dup: ", proto_name)
        sys.exit(-1)
    p = Proto(proto_name, p_list)
    proto_dict[proto_name] = p
    proto_list.append(p)

def parse_proto(proto_buf):
    L = [line.strip() for line in proto_buf.split('\n')]
    L = [line.split('#')[0].strip() for line in L if line and line[0] != '#']

    proto_dict = {}

    start_proto = False
    p_list = []
    for line in L:
        if line.find('===') >= 0: # 一个proto结束
            if start_proto:
                new_proto(proto_dict, proto_name, p_list)
                p_list = []
                proto_name = ''
            continue
        idx = line.find('=')
        if idx > 0:
            start_proto = True
            proto_name = line[:idx]
        elif idx ==-1:
            # 成员
            key_words = line.split(' ')
            if len(key_words) < 2:
                print('proto error:', line)
                sys.exit(-1)
            p_list.append([k.strip() for k in key_words])
            continue
        else:
            print("Error, empty protocol name.")
            sys.exit(-1)

    return proto_dict

def gen_go_proto(proto_dict):
    f = open(os.path.join('./', 'proto.go'), 'w')
    f.write("""package protocol
import (
. "misc/packet"
)
""")
    f.write("\ntype PKT_null struct {}\n\n")
    f.write("""func Decode_null(reader *Packet) (tbl PKT_null, err error) {\n	return\n}\n\n""")
    for k, p in proto_dict.items():
        f.write("""type PKT_%s struct {\n""" % (k))
        for field in p.fields:
            ary_buf = ""
            if field.array: ary_buf = "[]"
            f.write("\tF_%s %s%s\n" % (field.name, ary_buf, field.go_type()))
        f.write("}\n\n")

    for k, p in proto_dict.items():
        f.write("""func Decode_%s(reader *Packet) (tbl PKT_%s, err error) {\n"""
                % (k, k))
        for field in p.fields:
            if field.is_basic_type():
                func_str = "reader.%s()" % (field.go_type_func())
            else:
                func_str = "Decode_%s(reader)" % (field.type)

            if field.array:
                f.write("""\t{
		narr, err := reader.ReadU16()
		checkErr(err)

		tbl.F_%s = make([]%s, narr)
		for i := 0; i < int(narr); i++ {
			tbl.F_%s[i], err = %s
			checkErr(err)
		}
	}\n\n""" % (field.name, field.go_type(), field.name, func_str))
            else:
                f.write("\ttbl.F_%s, err = %s\n\tcheckErr(err)\n" \
                            % (field.name, func_str))
        f.write("\treturn\n}\n\n")

    f.close()


def gen_go_api(api_dict):
    f = open(os.path.join('./', 'api.go'), 'w')
    f.write("""package protocol\n\n
const(\n""")
    for k, api in api_dict.items():
        f.write("%s=%s // %s\n" % (str.upper(api.name), api.type, api.desc))
    f.write(")\n\n")

    f.write("var RCode = map[int16]string{\n")
    for k, api in api_dict.items():
        f.write("\t%s: \"%s\", // %s\n" % (api.type, api.name, api.desc))
    f.write("}\n\n")

    f.write('''type ApiInfo struct{
    Api int16
    ApiName string
    Payload interface{}
    Desc string
    Type int8
}
    var ApiMap = map[string]*ApiInfo{
''')
    for k, api in api_dict.items():
        f.write("\t\"%s\": &ApiInfo{%s,\"%s\",&PKT_%s{},\"%s\",%s},\n" % (api.name, api.type, api.name, api.payload,api.desc,api.flag))
    f.write("}\n\n")

    f.close()

    f = open(os.path.join('./', 'dispatcher.go'), 'w')

    f.write("""package gs
    import (
    . "protocol"
	"base/tlog"
    )
    
    func Handle(api int16, payload interface{}) {
    switch api {
    """)

    for k, api in api_dict.items():
        if api.type >= 20000: continue
        if len(api.name) < 4 or api.name[-4:] != "_req": continue
        f.write("""\tcase %s: P_%s(payload.(PKT_%s))\n""" %
                (api.type, api.name, api.payload))

    f.write("""default:
    tlog.Log("Unknown msg:", api)
}
}
""")
    f.close()

    f = open(os.path.join('./', 'decode.go'), 'w')
    f.write("""package protocol
import (
. "misc/packet"
)

func DecodeMsg(p *[]byte) (int16, interface{}) {
reader := Reader(*p)
api, err := reader.ReadS16()
if err != nil {
return -1, nil
}
var msg interface{}
switch api {
""")
    for k, api in api_dict.items():
        if len(api.name) < 4: continue
        f.write("""\tcase %s: msg, _ = Decode_%s(reader)\n""" % (api.type, api.payload))
    f.write("""default: return -1, nil
}
return api, msg
}
""")
    f.close()

def gen_go_err(err_dict):
    f = open(os.path.join('./', 'errcode.go'), 'w')
    f.write("""package protocol\n\nconst (\n""")
    for _, v in err_dict.items():
        code, name, desc = v
        f.write('''%s = %d // %s\n''' % (name, code, desc))
    f.write(')\n')

    f.write("\nvar ErrMap = map[int16]string{\n")
    for _, v in err_dict.items():
        code, name, desc = v
        f.write('''%s : \"%s\",\n''' % (code, desc))
    f.write('}\n')
    f.close()

def gen_go_const(const_list):
    f = open(os.path.join('./', 'const.go'), 'w')
    f.write("""package protocol\n\nconst (\n""")
    for item in const_list:
        if item[0] == '#':
            f.write("\n// %s\n" % (item[1:]))
            continue
        code, name, desc = item
        if type(code) == int:
            f.write('''%s = %s // %s\n''' % (name, code, desc))
        else:
            f.write('''%s = "%s" // %s\n''' % (name, code, desc))
    f.write(')\n')
    f.close()

def parse_errcode(errcode_buf):
    lines = [i.strip() for i in errcode_buf.split('\n')]
    lines = [l for l in lines if (l and l[0] != '#')]
    d = {}
    for line in lines:
        if not line: continue
        el = line.split('-')
        if len(el) != 3:
            print('Error errcode:', el)
            continue
        code, name, desc = el
        code = int(code)
        name = name.upper()
        if code in d:
            print('Dup errcode', el)
            continue
        d[code] = (code, name, desc)
    return d

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

        try:
            code = int(code)
        except ValueError:
            pass
        name = name.upper()
        d.append((code, name, desc))
    return d


# 检查无用协议是否被用到。
def check_protocol(proto_dict, api_dict):
    payload_dict = {}
    for _, api in api_dict.items():
        payload_dict[api.payload] = True

    for pdname in proto_dict:
        if pdname not in payload_dict:
            used = False
            for _, pd in proto_dict.items():
                if pd.name == pdname:
                    continue
                for field in pd.fields:
                    if field.type == pdname:
                        used = True
                        break
                if used: break
            if not used:
                print('protocol %s unused!' % pdname)
        

def parse(api_buf1, api_buf2, proto_buf, errcode_buf, const_buf):
    api_dict = parse_api(api_buf1,1)
    api_dict.update(parse_api(api_buf2,2))
    proto_dict = parse_proto(proto_buf)
    errcode_dict = parse_errcode(errcode_buf)
    const_list = parse_const(const_buf)

    check_protocol(proto_dict, api_dict)

    gen_go_proto(proto_dict)
    gen_go_api(api_dict)
    gen_go_err(errcode_dict)
    gen_go_const(const_list)

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('usage: ./parse.py proto_dir [gen_dir]')
        sys.exit(0)

    path_pre = sys.argv[1]
    try:
        api_buf1 = open(os.path.join(path_pre, 'api.txt'), 'r').read()
        proto_buf1 = open(os.path.join(path_pre, 'proto.txt'), 'r').read()
        api_buf2 = open(os.path.join(path_pre, 'intra_api.txt'), 'r').read()
        proto_buf2 = open(os.path.join(path_pre, 'intra_proto.txt'), 'r').read()

        proto_buf = proto_buf1 + '\n' + proto_buf2
        errcode_buf = open(os.path.join(path_pre, 'error.txt'), 'r').read()
        const_buf = open(os.path.join(path_pre, 'const.txt'), 'r').read()
    except (IOError, e):
        print('Open proto file failed:', e)
        sys.exit(0)

    parse(api_buf1,api_buf2, proto_buf, errcode_buf, const_buf)
