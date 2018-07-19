from optparse import OptionParser
from gw.gwservice import GWService
from webclient.webclient import WebClient
from constants import ACTION, ENV

parser = OptionParser()
parser.add_option('-e', '--env', dest='env', default='dev',
                  help='ENV should be [dev|prod]',
                  metavar='ENV')
parser.add_option('-t', '--action', dest='action', default='run',
                  help='ACTION should be [build|run|stop]',
                  metavar="ACTION")

(options, args) = parser.parse_args()


def run():
    services = [
        GWService(options.env),
        WebClient(options.env)
    ]

    if options.action == ACTION.BUILD:
        for s in services:
            if s.build() is False:
                print "fail to build: ", s.name()
                return

    elif options.action == ACTION.RUN:
        for s in services:
            if s.build() is False:
                print "fail to build: ", s.name()
                return

            if s.deploy() is False:
                print "fail to deploy: ", s.name()
                return

            if s.run() is False:
                print "fail to run: ", s.name()
                return

    elif options.action == ACTION.STOP:
        for s in services:
            s.stop()


if __name__ == '__main__':
    run()
