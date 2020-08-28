import subprocess
import json
import os.path
import sys
import urllib.parse

running = True
url = 'http://localhost:8080'
tks = []

class CmdException(Exception):
    def __init__(self, cause):
        self.cause = cause

    def str(self):
        return self.cause

def exec_cmd(cmd, expect_reply=True):
    res = subprocess.run(cmd, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if res.returncode != 0:
        print(res.stderr)
        sys.exit(res.returncode)

    if expect_reply and len(res.stdout) == 0:
        raise CmdException('Got an empty reply')
    elif not expect_reply and len(res.stdout) > 0:
        raise CmdException('Unexpected reply: ' + res.stdout.decode('utf-8'))

    if not expect_reply:
        return None
    else:
        try:
            return json.loads(res.stdout)
        except:
            raise CmdException(res.stdout.decode('utf-8'))

def get(url):
    cmd = [
        'curl',
        '-H',
        'Content-Type: application/json',
        url,
    ]
    return exec_cmd(cmd, expect_reply=True)

def post(url):
    cmd = [
        'curl',
        '-X',
        'POST',
        '-H',
        'Content-Type: application/json',
        url,
    ]
    exec_cmd(cmd, expect_reply=False)

def normalize_addr(addr):
    while addr[-1] == '/':
        addr = addr[:-1]
    return addr

def new_token(host, game):
    game = urllib.parse.quote(game, safe='')
    url = normalize_addr(host) + '/' + os.path.join('run', 'new', game)
    tk = get(url)
    global tks
    tks.append(tk['Token'])

def get_time(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', 'timer', token)
        t = get(url)
        print('Token: {} - Time: {}'.format(token, t['Time']))
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def get_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', 'splits', token)
        s = get(url)
        print(json.dumps(s, indent=4, sort_keys=False))
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def reset_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'reset')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def start_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'start')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def adv_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'split')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def undo_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'undo')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def skip_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'skip')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def pausetoggle_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'pause-toggle')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)

def save_splits(host, tk_idx):
    try:
        tk_idx = int(tk_idx)-1
        token = tks[tk_idx]
        url = normalize_addr(host) + '/' + os.path.join('run', token, 'save')
        post(url)
        print('OK!')
    except:
        print('Invalid token index "{}"'.format(tk_idx))
        list_tk(None)


#tk = new_token(url, 'JJAT%20(any%25)')
#t = get_time(url, tk)
#sp = get_splits(url, tk)
#
#print('token: {}'.format(tk))
#print('time: {}'.format(t))
#print('splits: {}'.format(repr(sp)))

def on_quit(_):
    global running
    running = False

def set_url(_, new_url):
    global url
    url = new_url

def get_url(_):
    print('Current host: {}'.format(url))

def list_tk(_):
    print('  Tokens:')
    for i in range(len(tks)):
        print('    {} - {}'.format(i+1, tks[i]))

def main():
    opts = [
        ('List tokens', (list_tk,)),
        ('New token', ('Insert the game/category: ', new_token,)),
        ('Get run time', ('Select a token (from "List tokens"): ', get_time,)),
        ('Retrieve the run\'s splits', ('Select a token (from "List tokens"): ', get_splits,)),
        ('Reset the current run', ('Select a token (from "List tokens"): ', reset_splits,)),
        ('Start a new run', ('Select a token (from "List tokens"): ', start_splits,)),
        ('Split to the next segment', ('Select a token (from "List tokens"): ', adv_splits,)),
        ('Undo the latest segment', ('Select a token (from "List tokens"): ', undo_splits,)),
        ('Skip the current segment', ('Select a token (from "List tokens"): ', skip_splits,)),
        ('Pause-toggle the run\'s timer', ('Select a token (from "List tokens"): ', pausetoggle_splits,)),
        ('Save the complete run', ('Select a token (from "List tokens"): ', save_splits,)),
        ('Print host', (get_url,)),
        ('Set host', ('Insert the runs host/url: ', set_url,)),
        ('Quit', (on_quit,)),
    ]

    print('Select an option:\n')
    for i in range(len(opts)):
        print('{} - {}'.format(i+1, opts[i][0]))

    i = input('\n  > ').strip()
    print('')
    try:
        i = int(i)-1
        cmd = opts[i][1]
        f = cmd[-1]
        if len(cmd) > 1:
            arg = input(cmd[0])
            f(url, arg)
        else:
            f(url)
    except e as Exception:
        print(e)
        print('Invalid option "{}"'.format(i))
    print('')

if __name__ == '__main__':
    while running:
        main()
