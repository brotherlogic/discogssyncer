import os
import subprocess

# Update to the latest version
for line in os.popen('git fetch -p -q; git merge -q origin/master').readlines():
    print line.strip()
    
# Move the old version over
for line in os.popen('cp discogssyncer oldsync').readlines():
    print line.strip()

# Rebuild
for line in os.popen('go build ./...').readlines():
    print line.strip()

# Rebuild
for line in os.popen('go build').readlines():
    print line.strip()

size_1 = os.path.getsize('./oldsync')
size_2 = os.path.getsize('./discogssyncer')

if size_1 != size_2:
    for line in os.popen('killall discogssyncer').readlines():
        pass
    subprocess.Popen(['./discogssyncer', '--port' ,'50051'])
