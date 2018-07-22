import yaml
from constants import ENV


class ServiceConfig:
    @staticmethod
    def dataDirName():
        # it must be 'data' which is agreement with build.yaml
        return 'data'

    def __init__(self, c, env):
        self.c = c
        self.env = env

    def name(self):
        """
        query the service name which will
        be used to create service.
        they must be same as the module name
        """
        return self.c['name']

    def bins(self):
        """
        bins includes all the binary executable files
        """
        if self.c['bin'] is None:
            return []
        return self.c['bin'][self.env]

    def data(self):
        """
        data files includes the env config files
        """
        if self.c['data'] is None:
            return []
        return self.c['data'][self.env]

    def buildCmd(self):
        """
        bash command for building the service
        """
        if self.c['build'] is None:
            return None
        return self.c['build'][self.env]

    def runCmd(self):
        """
        bash command for running the service
        """
        if self.c['run'] is None:
            return None
        return self.c['run'][self.env]

    def path(self):
        """
        service src code root directory
        """
        return self.c['path']


class Config:
    def __init__(self, fp, env):
        self.fp = fp
        self.c = None
        self.env = env

    def load(self):
        with open(self.fp, 'r') as stream:
            self.c = yaml.load(stream)

    def serviceConfigs(self):
        return map(lambda s: ServiceConfig(s, self.env), self.c.values())
