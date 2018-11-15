import os
from os import path
import time
from datetime import datetime

from executor import execute

CURRENT_VERSION_COMMAND = 'git rev-parse --short HEAD'
REMOTE_VERSION_COMMAND = "git ls-remote | grep master | awk '{print $1}'"
UPDATE_COMMAND = 'git pull'
RUN_SERVICE_COMMAND = 'sh ./prodrun.sh'
UPDATE_INTERVAL = 3 * 60  # 3 minutes


def log(msg):
    ct = datetime.now()
    print "t={}, msg={}".format(ct, msg)


def run_once(idx):
    """
    Follow the procedure
    1. record current commit id
    2. git pull the repo
    3. compare the pulled commit id with recorded commit id & done for this loop if they the same
    4. build the whole system and run it if build successfully
    """
    log("start ci#", idx)
    current_version = execute(CURRENT_VERSION_COMMAND, capture=True)
    latest_version = execute(REMOTE_VERSION_COMMAND, capture=True)

    if latest_version.startswith(current_version):
        print "no need to update code"
        return

    log("need update code, updating...")
    execute(UPDATE_COMMAND)
    log("finish update code, starting service...")
    execute(RUN_SERVICE_COMMAND)
    log("finish starting service")


def run():
    update_counter = 0
    while(True):
        try:
            run_once(update_counter)
        except Exception as e:
            log("fail to run_once, err: {}".format(e))
        finally:
            update_counter += 1
            time.sleep(UPDATE_INTERVAL)


if __name__ == '__main__':
    run()
