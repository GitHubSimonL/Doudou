#!/usr/bin/env python
#coding=utf8

"""
GRPC通讯中的数据类型与协议类型自动转换器
"""
import sys, os
import parse

servlName = ""
servuName = ""
basicTypeConvert = False

class MessageParam(object):
	def __init__(self, t, n, r, sn):
		super(MessageParam, self).__init__()
		self.type = t
		self.name = n
		self.repeated = r
		self.struct_name = sn

	def getMsgType(self, ptd):
		sn = self.struct_name.replace("PKT_", "")
		proto = ptd[sn]
		for field in proto.fields:
			if field.name == self.name:
				return field.go_type()


	def writeType2Msg(self, f, ptd):
		msgname = ""
		names = self.name.split('_')
		for i in range(0, len(names)):
			msgname += names[i].capitalize()
		t = self.getMsgType(ptd)
		if t == "int8":
			if self.repeated == True:
				f.write("""
		F_%s: int32_to_int8_array(p.%s),""" % (self.name, msgname))
			else:
				f.write("""
		F_%s: int8(p.%s),""" % (self.name, msgname))
		elif t == "int16":
			if self.repeated == True:
				f.write("""
		F_%s: int32_to_int16_array(p.%s),""" % (self.name, msgname))
			else:
				f.write("""
		F_%s: int16(p.%s),""" % (self.name, msgname))
		elif t == "PKT_rawdata":
			f.write("""
		F_%s: p.%s,""" % (self.name, msgname))
		elif t in ['int32', 'int64', 'string', 'bool']:
			f.write("""
		F_%s: p.%s,""" % (self.name, msgname))
		else:
			sn = ""
			sns = t.replace('PKT_', '').split('_')
			for i in range(0, len(sns)):
				sn += sns[i].capitalize()
			if self.repeated == True:
				f.write("""
		F_%s: %s_%s2MsgArray(p.%s),""" % (self.name, servuName, sn, msgname))
			else:
				f.write("""
		F_%s: %s_%s2Msg(*p.%s),""" % (self.name, servuName, sn, msgname))
			

	def writeType2MsgArray(self, f, ptd):
		pass

	def writeMsg2Type(self, f, ptd):
		msgname = ""
		names = self.name.split('_')
		for i in range(0, len(names)):
			msgname += names[i].capitalize()
		t = self.getMsgType(ptd)
		if t == "int8":
			if self.repeated == True:
				f.write("""
		%s: int8_to_int32_array(m.F_%s),""" % (msgname, self.name))
			else:
				f.write("""
		%s: int32(m.F_%s),""" % (msgname, self.name))
		elif t == "int16":
			if self.repeated == True:
				f.write("""
		%s: int16_to_int32_array(m.F_%s),""" % (msgname, self.name))
			else:
				f.write("""
		%s: int32(m.F_%s),""" % (msgname, self.name))
		elif t == "PKT_rawdata":
			f.write("""
		%s: m.F_%s,""" % (msgname, self.name))
		elif t in ['int32', 'int64', 'string', 'bool']:
			f.write("""
		%s: m.F_%s,""" % (msgname, self.name))
		else:
			sn = ""
			sns = t.replace('PKT_', '').split('_')
			for i in range(0, len(sns)):
				sn += sns[i].capitalize()
			if self.repeated == True:
				f.write("""
		%s: %s_Msg2%sArray(m.F_%s),""" % (msgname, servuName, sn, self.name))
			else:
				f.write("""
		%s: &%s,""" % (msgname, self.name))

	def writeMsg2TypeArray(self, f, ptd):
		pass
		


class Message(object):
	def __init__(self, n, ctn):
		super(Message, self).__init__()
		self.name = n
		self.convert_to_name = ctn
		self.params = {}
		self.struct_params = {}

	def addParam(self, t, n, r, sn, ptd):
		param = MessageParam(t, n, r, sn)
		self.params[n] = param
		t = param.getMsgType(ptd)
		if (not t in ['int8', 'int16', 'int32', 'int64', 'string', 'bool', 'PKT_rawdata']) and r == False:
			sn = ""
			sns = t.replace('PKT_', '').split('_')
			for i in range(0, len(sns)):
				sn += sns[i].capitalize()
			self.struct_params[n] = sn

	def writeFuncHeadType2Msg(self, f):
		f.write("""
func %s_%s2Msg(p %s) %s {
	return %s{""" % (servuName, self.name, self.name, self.convert_to_name, self.convert_to_name))

	def writeFuncHeadType2MsgArray(self, f):
		f.write("""
func %s_%s2MsgArray(p []*%s) []%s {
	infos := make([]%s, 0)
	for _, item := range p {
		infos = append(infos, %s_%s2Msg(*item))
	}
	return infos
}
""" % (servuName, self.name, self.name, self.convert_to_name, self.convert_to_name, servuName, self.name))

	def writeFuncHeadMsg2Type(self, f):
		f.write("""
func %s_Msg2%s(m %s) %s {""" % (servuName, self.name, self.convert_to_name, self.name))
		for k, sp in self.struct_params.items():
			f.write("""
	%s := %s_Msg2%s(m.F_%s)""" % (k, servuName, sp, k))
		f.write("""
	return %s{""" % (self.name))

	def writeFuncHeadMsg2TypeArray(self, f):
		f.write("""
func %s_Msg2%sArray(m []%s) []*%s {
	infos := make([]*%s, 0)
	for _, item := range m {
		info := %s_Msg2%s(item)
		infos = append(infos, &info)
	}
	return infos
}
""" % (servuName, self.name, self.convert_to_name, self.name, self.name, servuName, self.name))

	def writeFuncTailType2Msg(self, f):
		f.write("""
	}
}
""")

	def writeFuncTailType2MsgArray(self, f):
		pass

	def writeFuncTailMsg2Type(self, f):
		f.write("""
	}
}
""")

	def writeFuncTailMsg2TypeArray(self, f):
		pass



def parse_proto(proto_buf, btd):
	tokens = proto_buf.replace("convert:", " typeconvert ").replace("{", " mstart ").replace("}", " mend ").replace(";", " ; ").split()
	tlen = len(tokens)
	proto_dict = {}

	mstart = False
	name = ""
	cname = ""
	pt = ""
	pn = ""
	pr = False

	for i in range(0, tlen):
		if mstart == False and tokens[i] == "typeconvert":
			mstart = True
			name = ""
			cname = ""
			pt = ""
			pn = ""
			pr = False
		elif mstart == True and tokens[i] == "mend":
			mstart = False
			name = ""
			cname = ""
			pt = ""
			pn = ""
			pr = False
		elif mstart == True and (tokens[i] == "message" or tokens[i] == "mstart"):
			continue
		elif mstart == True and cname == "":
			cname = tokens[i]
		elif mstart == True and name == "":
			name = tokens[i]
			message = Message(name, cname)
			proto_dict[name] = message
		elif mstart == True and name != "" and tokens[i] == "repeated":
			pr = True
		elif mstart == True and name != "" and pt == "":
			pt = tokens[i]
		elif mstart == True and name != "" and pt != "" and pn == "":
			pn = tokens[i]
			proto_dict[name].addParam(pt, pn, pr, cname, btd)
		elif mstart == True and name != "" and pt != "" and pn != "" and tokens[i] == ";":
			pr = False
			pt = ""
			pn = ""
		
	return proto_dict


def gen_go_proto(proto_txt_dict, proto_dict):
	if len(proto_dict) <= 0:
		return
	f = open(os.path.join('./', '%s_type_convert.go' % servlName), 'w')
	f.write("""
package protobuf
import (
	. "protocol"
)
""")
	global basicTypeConvert
	if basicTypeConvert == False:   # 确保全局范围只生成一次如下基础类型转换代码
		f.write("""
func int8_to_int32_array(ns []int8) []int32 {
	rets := make([]int32, 0)
	for _,n := range ns {
		rets = append(rets, int32(n))
	}
	return rets
}

func int16_to_int32_array(ns []int16) []int32 {
	rets := make([]int32, 0)
	for _,n := range ns {
		rets = append(rets, int32(n))
	}
	return rets
}

func int32_to_int8_array(ns []int32) []int8 {
	rets := make([]int8, 0)
	for _,n := range ns {
		rets = append(rets, int8(n))
	}
	return rets
}

func int32_to_int16_array(ns []int32) []int16 {
	rets := make([]int16, 0)
	for _,n := range ns {
		rets = append(rets, int16(n))
	}
	return rets
}
""")
		basicTypeConvert = True
	for _, proto in proto_dict.items():
		proto.writeFuncHeadType2Msg(f)
		for _, param in proto.params.items():
			param.writeType2Msg(f, proto_txt_dict)
		proto.writeFuncTailType2Msg(f)

		proto.writeFuncHeadType2MsgArray(f)
		for _, param in proto.params.items():
			param.writeType2MsgArray(f, proto_txt_dict)
		proto.writeFuncTailType2MsgArray(f)

		proto.writeFuncHeadMsg2Type(f)
		for _, param in proto.params.items():
			param.writeMsg2Type(f, proto_txt_dict)
		proto.writeFuncTailMsg2Type(f)

		proto.writeFuncHeadMsg2TypeArray(f)
		for _, param in proto.params.items():
			param.writeMsg2TypeArray(f, proto_txt_dict)
		proto.writeFuncTailMsg2TypeArray(f)

	f.close()

def parse_type_convert(proto_txt_buf, proto_buf):
	proto_txt_dict = parse.parse_proto(proto_txt_buf)
	proto_dict = parse_proto(proto_buf, proto_txt_dict)

	gen_go_proto(proto_txt_dict, proto_dict)

if __name__ == "__main__":
	if len(sys.argv) < 2:
		print('usage: ./parse_type_convert.py proto_dir')
		sys.exit(0)

	path_pre = sys.argv[1]

	proto_txt_buf_1 = open(os.path.join(path_pre, 'proto.txt'), 'r').read() 
	proto_txt_buf_2 = open(os.path.join(path_pre, 'intra_proto.txt'), 'r').read()
	proto_txt_buf = proto_txt_buf_1 + proto_txt_buf_2

	for _,_,fs in os.walk(path_pre):
		for f in fs:
			if os.path.splitext(f)[1] == '.proto':
				servlName = f.replace(".proto", "")
				servuName = servlName.capitalize()
				proto_buf = open(os.path.join(path_pre, f), 'r').read()
				parse_type_convert(proto_txt_buf, proto_buf)
