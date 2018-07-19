import os
import random
from executor import execute

from constants import GoRootDir
from common.deployer import ServiceDeployer
from common.proc import RerunProc, KillProc


class GWService(ServiceDeployer):
    @staticmethod
    def name():
        return 'gateway'

    def __init__(self, mode):
        super(GWService, self).__init__(mode)

    def build(self):
        GoRootDir()
        execute('mkdir -p "{}"'.format(self.buildDir()))
        execute(
            'go build  -o "{}/main" cmd/gw/main.go'.format(self.buildDir()))
        return True

    def deploy(self):
        GoRootDir()
        execute('mkdir -p "{}"'.format(self.workDir()))
        execute('mv "{}/main" "{}/main"'.format(self.buildDir(), self.workDir()))
        return True

    def run(self):
        GoRootDir()

        pidFilePath = '{}/pid'.format(self.workDir())
        command = '{}/main'.format(self.workDir())
        return RerunProc(pidFilePath, command)

    def stop(self):
        pidFilePath = '{}/pid'.format(self.workDir())
        KillProc(pidFilePath)
