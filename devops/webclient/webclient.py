from executor import execute
from constants import GoRootDir
from common.deployer import ServiceDeployer
from common.proc import RerunProc, KillProc


class WebClient(ServiceDeployer):
    @staticmethod
    def name():
        return "webclient"

    def __init__(self, mode):
        super(WebClient, self).__init__(mode)

    def build(self):
        GoRootDir()
        execute('mkdir -p "{}"'.format(self.buildDir()))

        return execute("cp -r demo/* {}".format(self.buildDir()))

    def deploy(self):
        GoRootDir()
        execute('mkdir -p "{}"'.format(self.workDir()))
        return execute('cp -r {}/* {}/'.format(self.buildDir(), self.workDir()))

    def run(self):
        GoRootDir()

        pidFilePath = '{}/pid'.format(self.workDir())
        command = 'swank --path {}'.format(self.workDir())
        return RerunProc(pidFilePath, command)

    def stop(self):
        pidFilePath = '{}/pid'.format(self.workDir())
        KillProc(pidFilePath)
