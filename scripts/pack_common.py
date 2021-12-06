#!/usr/bin/env python
# -*- coding: utf-8 -*-

#import zipfile
import os
import sys
import ConfigParser
import shutil


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


def pack_l3(includeTool):
    config = ConfigParser.ConfigParser()
    config.read("./config.ini.template")
    #gate
    deb_dir = "./data/deb/"
    deb_name = "l3"
    deb_path = deb_dir + deb_name
    deb_install_path = deb_path + "/data/l3/"
    if os.path.exists(deb_path):
        #os.rename(deb_path, deb_path + ".bak")
        shutil.rmtree(deb_path)

    # publish path
    try:
        os.makedirs("./data/publish")
    except:
        pass

    #supervisor config
    os.makedirs(deb_install_path)
    os.makedirs(deb_install_path + "log")
    shutil.copy("l3.conf", deb_install_path)

    os.system("git log | head -1 >>" + deb_install_path + "/MD5")

    shutil.copy("./kill.sh", current_dir)

    current_sections = {}.fromkeys(("GLOBAL", "LOG", "GATE", "GS_TEMP", "CSV"))
    current_sections["LOG"] = ("log_maxsize", "default_path")
    current_sections["GATE"] = ("intra_ip", "intra_port")
    current_sections["GS_TEMP"] = ("server_id", "mongo_user", "mongo_pass",
                                   "mongo_addr", "mongo_dbname")
    set_config(config, current_sections, current_dir + "/config.ini.template")

    #gate
    current_dir = deb_install_path + "/gate"
    current_sections = {}.fromkeys(
        ("GLOBAL", "LOG", "MSG", "CSV", "GATE", "PLATFORM", "ACCOUNT",
         "GM_CMD", "MONGODB_DUMP", "MSG_STAT"))
    os.makedirs(current_dir)
    os.makedirs(current_dir + "/log")
    os.makedirs(current_dir + "/gd_config")
    shutil.copy("bin/gate_main", current_dir)
    set_config(config, current_sections, current_dir + "/config.ini.template")
    copy_dir("gd_config", current_dir + "/gd_config")

    #gs
    current_dir = deb_install_path + "/gs"
    current_sections = {}.fromkeys(
        ("GLOBAL", "LOG", "MSG", "CSV", "PLATFORM", "", "GATE", "CENTER",
         "CHAT", "MSG_STAT", "GS_TEMP", "GM_CMD", "MONGODB_DUMP"))
    current_sections["GATE"] = ("intra_ip", "intra_port")
    current_sections["CENTER"] = ("ip", "port")
    os.makedirs(current_dir)
    os.makedirs(current_dir + "/log")
    os.makedirs(current_dir + "/gd_config")
    os.makedirs(current_dir + "/script")
    shutil.copy("bin/gs_main", current_dir)
    set_config(config, current_sections, current_dir + "/config.ini.template")
    copy_dir("gd_config", current_dir + "/gd_config")
    copy_dir("script", current_dir + "/script")

    #backup_mongodb
    current_dir = deb_install_path + "/backup_mongodb"
    os.makedirs(current_dir)
    os.makedirs(current_dir + "/log")
    copy_dir("src/scripts/backup_mongodb", current_dir)

    os.makedirs(deb_path + "/DEBIAN/")

    f = open(deb_path + "/DEBIAN/control", "w")

    f.write('''Source: l3
Section: unknown
Priority: extra
Version: 1
Maintainer: l3 <l3@igg.com>
Homepage: <insert the upstream URL, if relevant>
Package: l3
Architecture: amd64
Depends:
Description:l3
''')
    f.close()

    os.system("dpkg-deb --build " + deb_path)

    os.rename(deb_path + ".deb", "./data/publish/" + deb_name + ".deb")

    print "build l3 ok"


def pack_data():
    config = ConfigParser.ConfigParser()
    config.read("./config.ini.template")
    #gate
    deb_dir = "./data/deb/config/"
    deb_name = "l3"
    deb_final_name = "l3_data"
    deb_path = os.path.join(deb_dir, deb_name)
    deb_install_path = os.path.join(deb_path, "./")
    if os.path.exists(deb_path):
        #os.rename(deb_path, deb_path + ".bak")
        shutil.rmtree(deb_path)

    #
    try:
        os.makedirs("./data/publish")
    except:
        pass

    os.makedirs(deb_install_path)
    os.system("git log | head -1 >>" + deb_install_path + "/MD5")

    #gate
    current_dir = deb_install_path + "/gate"
    os.makedirs(current_dir)
    os.makedirs(current_dir + "/gd_config")
    copy_dir("gd_config", current_dir + "/gd_config")

    #gs
    current_dir = deb_install_path + "/gs"
    os.makedirs(current_dir)
    os.makedirs(current_dir + "/gd_config")
    os.makedirs(current_dir + "/script")
    copy_dir("gd_config", current_dir + "/gd_config")
    copy_dir("script", current_dir + "/script")

    os.makedirs(deb_path + "/DEBIAN/")

    f = open(deb_path + "/DEBIAN/control", "w")

    f.write('''Source: l3
Section: unknown
Priority: extra
Version: 1
Maintainer: l3 <l3@igg.com>
Homepage: <insert the upstream URL, if relevant>
Package: l3
Architecture: amd64
Depends:
Description:l3
''')
    f.close()

    os.system("dpkg-deb -bZ bzip2 " + deb_path)

    os.rename(deb_path + ".deb", "./data/publish/" + deb_final_name + ".deb")

    print "build data ok"
