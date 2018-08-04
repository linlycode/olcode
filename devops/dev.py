from optparse import OptionParser

from constants import ACTION, ENV, ALL_SERVICES
from config import Config, ServiceConfig
from service import Service

import constants

parser = OptionParser()
parser.add_option('-e', '--env', dest='env', default=ENV.DEV,
                  help='ENV should be [dev|prod]',
                  metavar='ENV')
parser.add_option('-t', '--action', dest='action', default=ACTION.RUN,
                  help='ACTION should be [build|run|stop]',
                  metavar="ACTION")
parser.add_option('-s', '--service', dest='sname', default=ALL_SERVICES,
                  help='SERVICE should be [all|gw|demo|...]',
                  metavar='SERVICE')
parser.add_option('-p', '--workspace', dest='workspace', default=None,
                  help='WORKSPACE is the directory where all the service run',
                  metavar='WORKPSACE')
parser.add_option('-f', '--buildfile', dest='buildfile',
                  default='../build.yaml',
                  help='BUILDFILE is necessary for devops to build/running')

(options, args) = parser.parse_args()


def buildService(sc):
    """
    build service instance according to service name in sc
    """
    return Service(sc)


# TODO: move to a single file
CRED = '\33[31m'
CGREEN = '\33[32m'
CEND = '\33[91m'


def failPrint(s):
    print CRED, s, CEND


def succeedPrint(s):
    print CGREEN, s, CEND


def run():
    if options.workspace is not None:
        constants.WORK_DIR = options.workspace

    if not ENV.valid(options.env):
        print "invalid env=", options.env, ", must in ", ENV.allEnvs()
        return

    c = Config(options.buildfile, options.env)
    c.load()

    services = map(lambda sc: buildService(sc), c.serviceConfigs())

    if options.action == ACTION.BUILD:
        for s in services:
            if options.sname != ALL_SERVICES and options.sname != s.name():
                continue
            if s.build() is False:
                # TODO: prefer smarter way
                failPrint("(*1/2)fail to build:{}".format(s.name()))
                return
            else:
                succeedPrint("(1/2)succeed to build:{}".format(s.name()))

            if s.deploy() is False:
                failPrint("(*2/2)fail to deploy:{}".format(s.name()))
                return
            else:
                succeedPrint("(2/2)succeed to deploy:{}".format(s.name()))

    elif options.action == ACTION.RUN:
        for s in services:
            if options.sname != ALL_SERVICES and options.sname != s.name():
                continue
            if s.build() is False:
                failPrint("(*1/3)fail to build:{}".format(s.name()))
                return
            else:
                succeedPrint("(1/3)succeed to build:{}".format(s.name()))

            if s.deploy is False:
                failPrint("(*2/3)fail to deploy:{}".format(s.name()))
                return
            else:
                succeedPrint("(2/3)succeed to deploy:{}".format(s.name()))

            if s.run() is False:
                failPrint("(*3/3)fail to run:{}".format(s.name()))
                return
            else:
                succeedPrint("(3/3)succeed to run:{}".format(s.name()))

    elif options.action == ACTION.STOP:
        for s in services:
            if options.sname != ALL_SERVICES and options.sname != s.name():
                continue
            s.stop()


if __name__ == '__main__':
    run()
