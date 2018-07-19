import os
from executor import execute


def KillProc(pidFilePath):
    if os.path.exists(pidFilePath):
        with open(pidFilePath) as f:
            pid = f.readline()
            try:
                execute('kill -9 {}'.format(pid))
            except Exception:
                pass


def RerunProc(pidFilePath, command):
    # kill the last running process
    KillProc(pidFilePath)

    # we skip this error so that the existing process
    # will be killed next time
    execute('{} &'.format(command))

    # find the pid and save it for killing next time
    res = execute(
        'ps aux | grep "{}" | grep -v "grep"'.format(command), capture=True)

    pid = str(res).split()[1]
    with open(pidFilePath, 'w') as f:
        f.write(str(pid))

    return True
