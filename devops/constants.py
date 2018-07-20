import os
from os import path
from executor import execute


# env definition
class ENV:
    DEV = 'dev'
    PROD = 'prod'


# action definition
class ACTION:
    BUILD = 'build'
    RUN = 'run'
    STOP = 'stop'


BUILD_DIRNAME = ".build"
WORK_DIRNAME = ".work"


def GoRootDir():
    os.chdir(path.join(path.dirname(__file__), os.pardir))


execute("mkdir -p ./{}".format(BUILD_DIRNAME))
execute("mkdir -p ./{}".format(WORK_DIRNAME))


def GetBuildDir():
    return path.join(path.dirname(__file__), BUILD_DIRNAME)


def GetWorkDir():
    return path.join(path.dirname(__file__), WORK_DIRNAME)
