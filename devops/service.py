from abc import abstractmethod
from executor import execute

from common.proc import KillProc, RunProc
from constants import GetWorkDir, GoRootDir


class Service(object):
    def __init__(self, config):
        """
        :type config: config.ServiceConfig
        """
        self.c = config

    def name(self):
        return self.c.name()

    def build(self):
        GoRootDir()
        if self.c.buildCmd() is not None:
            execute(self.c.buildCmd())
        return True

    def deploy(self):
        GoRootDir()
        if not execute('mkdir -p {}'.format(self.workDir())):
            return False

        if not execute('rm -rf {0} && mkdir -p {0}'.format(self.workDir())):
            return False

        dataFiles = " ".join(self.c.data())
        if len(dataFiles) != 0:
            if not execute('mkdir -p {1}/{2} && cp -r {0} {1}/{2}'
                           .format(dataFiles, self.workDir(), self.c.dataDirName())):
                return False

        targetFiles = " ".join(self.c.bins())
        if len(targetFiles) != 0:
            if not execute('cp -r {} {}'.format(targetFiles, self.workDir())):
                return False

        execute('rm -rf {}'.format(targetFiles))
        return True

    def run(self):
        GoRootDir()
        if self.c.runCmd is not None:
            command = 'cd {} && {}'.format(self.workDir(), self.c.runCmd())
        return RunProc(self.workDir(), command)

    def stop(self):
        GoRootDir()
        KillProc(self.workDir())
        return

    def workDir(self):
        return '{}/{}'.format(GetWorkDir(), self.name())
