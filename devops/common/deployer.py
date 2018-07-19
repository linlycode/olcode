from abc import abstractmethod
from constants import GetBuildDir, GetWorkDir


class ServiceDeployer(object):
    @staticmethod
    @abstractmethod
    def name():
        pass

    def __init__(self, mode):
        self.mode = mode

    @abstractmethod
    def build(self):
        pass

    @abstractmethod
    def deploy(self):
        pass

    @abstractmethod
    def run(self):
        pass

    @abstractmethod
    def stop(self):
        pass

    def buildDir(self):
        return '{}/{}'.format(GetBuildDir(), self.name())

    def workDir(self):
        return '{}/{}'.format(GetWorkDir(), self.name())
