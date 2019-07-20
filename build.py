#!/usr/bin/env python3

import os
import subprocess
import shutil

OUT_DIR = "out"
GOPHER_PATH = "gopher.png"
MAZE_PATH = "maze"

if not os.path.isdir(OUT_DIR):
    os.mkdir(OUT_DIR)
else:
    shutil.rmtree(OUT_DIR)
    os.mkdir(OUT_DIR)

subprocess.run("go build -o {}/rome".format(OUT_DIR), shell=True)
subprocess.run("GOOS=windows GOARCH=amd64 go build -o {}/rome.exe".format(OUT_DIR), shell=True)
shutil.copy(GOPHER_PATH, OUT_DIR)
shutil.copytree(MAZE_PATH, "{}/maze".format(OUT_DIR))