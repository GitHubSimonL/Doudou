#!/usr/bin/env python
# -*- coding: utf-8 -*-

#import zipfile
import os
import sys
import ConfigParser
import shutil
import pack_common

def copy_dir(dirname, origdir):
    for dirpath, dirnames, filenames in os.walk(dirname):
        for filename in filenames:
            shutil.copy(os.path.join(dirpath,filename), os.path.join(origdir, filename))

def set_config(config, sections, orig_file):
    cf = ConfigParser.ConfigParser()
    for key,sc in sections.items():
        if config.has_section(key):
            cf.add_section(key)
            if sc == None :
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
    pack_common.pack_data()
