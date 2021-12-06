#!/usr/bin/env python
#coding=utf8
import os
import sys

class Func(object):
	"""docstring for Func"""
	def __init__(self, funcname, reqtype, replytype, flag):
		self.funcname = funcname
		self.reqtype = reqtype
		self.replytype= replytype
		self.flag = flag

class Service(object):
	"""docstring for Service"""
	def __init__(self, name):
		self.name = name
		self.func = []


def parse(buff):
	parse_protobuf(buff)

def parse_protobuf(buff):
	tokens = buff.replace("(", " ").replace(")", " ").replace("{", " ").replace("}", " ").split()
	i = 0
	service = None
	start = False
	end = False
	rpc = False
	funcname = ""
	reqtype = ""
	replytype = ""
	funcflag = True
	for i in range(0, len(tokens)):
		if start == False and tokens[i] != "service":
			continue
		elif start == False and tokens[i] == "service":
			start = True
		elif start == True and service == None:
			service = Service(tokens[i])
		elif end == False and tokens[i] == "rpc":
			rpc = True
		elif end == False and rpc == True and funcname == "":
			funcname = tokens[i]
		elif end == False and rpc == True and funcname != "" and reqtype == "":
			if tokens[i] == "stream":
				funcflag = False
				continue
			else:
				reqtype = tokens[i]
		elif end == False and rpc == True and funcname != "" and reqtype != "" and replytype == "":
			if tokens[i] == "returns":
				continue
			else:
				replytype = tokens[i]
				func = Func(funcname, reqtype, replytype, funcflag)
				rpc = False
				funcname = ""
				reqtype = ""
				replytype = ""
				funcflag = True
				service.func.append(func)
		elif end == False and tokens[i] == "message":
			end = True
		if end == True:
			break
	gen_go_proto(service)
			
def gen_go_proto(service):
	sname = service.name
	lsname = sname.lower()
	usname = sname.upper()
	f = open(os.path.join('./', '%s_grpc_api.go' % lsname), 'w')

	f.write("""package grpc_api\n
import (
    . "protobuf"
    "base/tlog"
    "helper"
    "errors"
    "misc"
    "time"
    "code.google.com/p/go.net/context"
    "google.golang.org/grpc"
)""")
	f.write("""
type %sFuncHandlers struct {""" % sname)
	for func in service.func:
		name = func.funcname
		reqtype = func.reqtype
		replytype = func.replytype
		flag = func.flag
		if flag == False:
			f.write("""
	func%s func(*%s, %s_%sClient, interface{})""" % (name, reqtype, sname, name))
		else:
			f.write("""
	func%s func(*%s, *%s, interface{})""" % (name, reqtype, replytype))
	f.write("""
}""")

	f.write("""
type %sGrpcCall struct {
    MsgType string
    Args interface{}
    Reply interface{}
    Extra interface{}
    Handle %sFuncHandlers
}	
type %sApi struct {
	addr string
    %s %sClient
    _%s_GRPC_MSG chan *%sGrpcCall
    _handle %sFuncHandlers
}
var (
	%sAPI %sApi
)
""" % (sname, sname, sname, usname, sname, usname, sname, sname, usname, sname))
	f.write("""
func (api *%sApi) MSGQ() chan *%sGrpcCall {
	return %sAPI._%s_GRPC_MSG
}
""" % (sname, sname, usname, usname))
	f.write("""
func init() {
	%sAPI = %sApi{}
	%sAPI._%s_GRPC_MSG = make(chan *%sGrpcCall, 102400)
	%sAPI._handle = %sFuncHandlers{}
}
""" % (usname, sname, usname, usname, sname, usname, sname))

	f.write("""
func Get%sClient(addr string) *%sApi {
	if addr == "" {
		tlog.LogErrf("rpc_addr not set.")
		return nil
	}
	api := %sApi{}
	api._%s_GRPC_MSG = make(chan *%sGrpcCall, 102400)
	api._handle = %sFuncHandlers{}
	api.addr = addr
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		tlog.LogErrf("Cannot connect to %s server,addr")
		return nil
	}
	api.%s = New%sClient(conn)
	tlog.Logf("%s server connected.")
	return &api
}
""" % (sname, sname, sname, usname, sname, sname, sname, usname, sname, sname))

	f.write("""
func Get%sMSGQ(api *%sApi) chan *%sGrpcCall{
	return api._%s_GRPC_MSG
}
""" % (sname, sname, sname, usname))

	f.write("""
func Dial%sServer(addr string) error {
	if addr == "" {
		tlog.LogErrf("rpc_addr not set.")
		return errors.New("rpc_addr not set")
	}
	%sAPI.addr = addr
    conn, err := grpc.Dial(addr, grpc.WithInsecure())
    if err != nil {
        tlog.LogErrf("Cannot connect to %s server, addr")
        return err
    }
    %sAPI.%s = New%sClient(conn)
    tlog.Logf("%s server connected.")
    return nil
}
func catch%sPanic(msg string) {
    if err := recover(); err != nil {
        helper.BackTrace(msg)
        tlog.LogErrf("Panic reconverd, err:", err)
    }
}
""" % (sname, usname, sname, usname, usname, sname, sname, sname))
	for func in service.func:
		name = func.funcname
		reqtype = func.reqtype
		replytype = func.replytype
		flag = func.flag
		if flag == False:
			f.write("""
func (api *%sApi) Register%s(fun func(*%s, %s_%sClient, interface{})) {
	api._handle.func%s = fun
}
""" % (sname, name, reqtype, sname, name, name))
		else:
			f.write("""
func (api *%sApi) Register%s(fun func(*%s, *%s, interface{})) {
	api._handle.func%s = fun
}
""" % (sname, name, reqtype, replytype, name))



		f.write("""
func (api *%sApi) Async%s(req %s, extra interface{}) {
	fun := func() {
		defer catch%sPanic("%s")
		if api.%s == nil {
			tlog.LogWarnf("%s is nil")
			return
		}""" % (sname, name, reqtype, sname, name, usname, usname))
		if flag == False:
			f.write("""
		reply, err := api.%s.%s(context.Background())
""" % (usname, name))
		else:
			f.write("""
		reply, err := api.%s.%s(context.Background(), &req)
""" % (usname, name))
		f.write("""
		if err != nil {
			tlog.LogWarnf("%s.%s error:", err)
			return
		}
		gc := &%sGrpcCall{
			MsgType: "%s",
			Args: &req,
			Reply: reply,
			Extra: extra,
			Handle: %sFuncHandlers{},
		}
		api._%s_GRPC_MSG <- gc
	}
	go fun()
}
""" % (usname, name, sname, name, sname, usname))


		if flag == False:
			f.write("""
func (api *%sApi) AsyncFun%s(req %s, extra interface{}, handle func(*%s, %s_%sClient, interface{})) {
""" % (sname, name, reqtype, reqtype, sname, name))
		else:
			f.write("""
func (api *%sApi) AsyncFun%s(req %s, extra interface{}, handle func(*%s, *%s, interface{})) {
""" % (sname, name, reqtype, reqtype, replytype))
		f.write("""
    fun := func() {
        defer catch%sPanic("%s")
        if api.%s == nil {
            tlog.LogWarnf("%s is nil")
            return
        }""" % (sname, name, usname, usname))
		if flag == False:
			f.write("""
        reply, err := api.%s.%s(context.Background())
""" % (usname, name))
		else:
			f.write("""
        reply, err := api.%s.%s(context.Background(), &req)
""" % (usname, name))
		f.write("""
        if err != nil {
            tlog.LogWarnf("%s.%s error:", err)
            return
        }
        gc := &%sGrpcCall{
            MsgType: "%s",
            Args: &req,
            Reply: reply,
            Extra: extra,
            Handle: %sFuncHandlers{func%s : handle},
        }
        api._%s_GRPC_MSG <- gc
    }
    go fun()
}
""" % (usname, name, sname, name,  sname, name, usname))


		if flag == False:
			f.write("""
func (api *%sApi) Sync%s(req %s, extra interface{}) (%s_%sClient, error) {""" % (sname, name, reqtype, sname, name))
		else:
			f.write("""
func (api *%sApi) Sync%s(req %s, extra interface{}) (*%s, error) {""" % (sname, name, reqtype, replytype))
		f.write("""
	defer catch%sPanic("%s")
	if api.%s == nil {
		tlog.LogWarnf("%s is nil")
		return nil,errors.New("%s is nil")
	}""" % (sname, name, usname, usname, usname))
		if flag == False:
			f.write("""
	reply, err := api.%s.%s(context.Background())""" % (usname, name))
		else:
			f.write("""
	reply, err := api.%s.%s(context.Background(), &req)""" % (usname, name))
		f.write("""
	return reply, err
}
""")
		if flag == True:
			f.write("""
func (api *%sApi) TTSync%s(req *%s, timeout time.Duration, extra interface{}) (*%s, error) {
	defer catch%sPanic("TTSync%s")
	if api.%s == nil {
		tlog.LogWarnf("%s is nil")
		return nil, errors.New("%s is nil")
	}

	r := misc.TTGRPC.Add()
	defer misc.TTGRPC.Remove(r.GetKey())

	go func(key string){
		defer catch%sPanic("TTSync%s")
		reply, err := api.%s.%s(context.Background(), req)
		if err != nil {
			tlog.LogErr("TTSync%s err: ", err)
			misc.TTGRPC.StopAsyncResult(key)
		}else{
			misc.TTGRPC.FillAsyncResult(key, reply)
		}
	}(r.GetKey())
	reply, err := r.GetResult(timeout)
	if err != nil {
		tlog.LogErr("TTSync%s err: ", err)
		return nil, err
	}else{
		return reply.(*%s), nil
	}
}""" % (sname, name, reqtype, replytype, sname, name, usname, usname, usname, sname, name, usname, name, name, name, replytype))

		f.write("""
func (api *%sApi) %sAck(req interface{}, reply interface{}, extra interface{}, handle %sFuncHandlers) {
	if api._handle.func%s == nil && handle.func%s == nil {
		tlog.LogWarnf("%s handle not registered or not set")
		return
	}
	Req := req.(*%s) """ % (sname, name, sname, name, name, name, reqtype))
		if flag == False:
			f.write("""
	Reply := reply.(%s_%sClient)""" % (sname, name))
		else:
			f.write("""
	Reply := reply.(*%s)""" % (replytype))
		f.write("""
	if api._handle.func%s != nil {
		api._handle.func%s(Req, Reply, extra)
	} else {
		handle.func%s(Req, Reply, extra)
	}
}
""" % (name, name, name))


	f.write("""
func (api *%sApi) Handle(msg *%sGrpcCall) {
	funcName := msg.MsgType
	switch funcName {
""" % (sname, sname))
	for func in service.func:
		name = func.funcname
		reqtype = func.reqtype
		replytype = func.replytype
		f.write("""
	case "%s":
		api.%sAck(msg.Args, msg.Reply, msg.Extra, msg.Handle)""" % (name, name))
	f.write("""
	}
}
""")

	f.close()



if __name__ == "__main__":
	if len(sys.argv) < 2:
		print('usage: ./parse_protobuf.py proto_dir')
		sys.exit(0)

	path_pre = sys.argv[1]
	grpc_path = "./src/grpc_api"
	if not os.path.exists(grpc_path):
		os.mkdir(grpc_path)
	for _,_,fs in os.walk(path_pre):
		for f in fs:
			if os.path.splitext(f)[1] == '.proto':
				buff = open(os.path.join(path_pre, f), 'r').read()
				parse(buff)

