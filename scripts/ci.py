import os
import time
from os import path

from executor import execute

CURRENT_VERSION_COMMAND = 'git rev-parse --short HEAD'
REMOTE_VERSION_COMMAND = "git ls-remote | grep master | awk '{print $1}'"
UPDATE_COMMAND = 'git pull'
RUN_SERVICE_COMMAND = 'sh ./prodrun.sh'
UPDATE_INTERVAL = 3 * 60 # 3 minutes

'''
follow the procedure
1. record current commit id
2. git pull the repo
3. compare the pulled commit id with recorded commit id & done for this loop if they the same
4. build the whole system and run it if build successfully
'''
def run_once(idx):
    print "start ci#", idx
    current_version = execute(CURRENT_VERSION_COMMAND, capture=True)
    latest_version = execute(REMOTE_VERSION_COMMAND, capture=True)

    if latest_version.startswith(current_version):
        print "no need to update code"
        return

    print "need update code, updating..."
    execute(UPDATE_COMMAND)
    print "finish update code, starting service..."
    execute(RUN_SERVICE_COMMAND)
    print "finish starting service"


def run():
    update_counter = 0
    while(True):
        try:
            run_once(update_counter)
        except Exception as e:
            print "fail to run_once, err:", e
        finally:
            time.sleep(UPDATE_INTERVAL)


if __name__ == '__main__':
    run()
