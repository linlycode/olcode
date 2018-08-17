import os
from executor import execute


def genPidFilePath(workDir):
    return os.path.join(workDir, 'pid')


def KillProc(workDir):
    pidFilePath = genPidFilePath(workDir)
    if os.path.exists(pidFilePath):
        with open(pidFilePath) as f:
            pid = f.readline()
            try:
                execute('kill -9 {}'.format(pid))
            except Exception:
                pass


def RunProc(workDir, command):
    # kill the last running process
    KillProc(workDir)

    # we skip this error so that the existing process
    # will be killed next time
    execute('{} &'.format(command))

    # find the pid and save it for killing next time
    # e.g., "mkdir gw && cd gw && ./gw -config ./data/dev.yaml &> gw.log"
    realCommand = command.split("&&")[-1].strip()
    logIndex = realCommand.find('&>')
    if logIndex > 0:
        realCommand = realCommand[:logIndex]

    res = execute(
        'ps aux | grep "{}" | grep -v "grep"'.format(realCommand), capture=True)

    try:
        pid = str(res).split()[1]
        int(pid)
    except ValueError:
        return False

    pidFilePath = genPidFilePath(workDir)
    with open(pidFilePath, 'w') as f:
        f.write(str(pid))

    return True
