#!/usr/bin/env python
# -*- coding: utf-8 -*-

#import zipfile
import os
import sys
import ConfigParser
import shutil
import pack_common

def copy_dir(dirname, origdir):
    for f in os.listdir(dirname):
        path = os.path.join(dirname, f)
        orig_path = os.path.join(origdir, f)
        if os.path.isdir(path):
            print "copy dir:", path, "to:", orig_path
            os.makedirs(orig_path)
            copy_dir(path, orig_path)
        if os.path.isfile(path):
            print "copy file:", path, "to:", orig_path
            shutil.copy(path, orig_path)


def set_config(config, sections, orig_file):
    cf = ConfigParser.ConfigParser()
    for key, sc in sections.items():
        if config.has_section(key):
            cf.add_section(key)
            if sc == None:
                for sub_key, value in config.items(key):
                    cf.set(key, sub_key, value)
            else:
                for sub_key in sc:
                    if config.has_option(key, sub_key):
                        cf.set(key, sub_key, config.get(key, sub_key))
    of = open(orig_file, "w")
    cf.write(of)
    of.close()


if __name__ == '__main__':
    include_tool = False
    if len(sys.argv) > 1:
        include_tool = (sys.argv[1] == 'true')
    pack_common.pack_l3(include_tool)
    pack_common.pack_data()
