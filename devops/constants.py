import os
from os import path
from executor import execute


# env definition
class ENV:
    DEV = 'dev'
    PROD = 'prod'

    @staticmethod
    def allEnvs():
        return [ENV.DEV, ENV.PROD]

    @staticmethod
    def valid(env):
        return env in ENV.allEnvs()


# action definition
class ACTION:
    BUILD = 'build'
    RUN = 'run'
    STOP = 'stop'


def GoRootDir():
    os.chdir(path.join(path.dirname(__file__), os.pardir))


WORK_DIR = ".work"
workDirCreated = False


def GetWorkDir():
    global workDirCreated
    wd = path.join(path.dirname(__file__), WORK_DIR)
    if not workDirCreated:
        execute("mkdir -p {}".format(wd))
        workDirCreated = True
    return wd
