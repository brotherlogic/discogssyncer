import os
import subprocess

current_hash = os.popen('git rev-parse HEAD').readlines()[0]
# Update to the latest version
for line in os.popen('go get -u github.com/brotherlogic/discogssyncer').readlines():
    print line.strip()
new_hash = os.popen('git rev-parse HEAD').readlines()[0]

    
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

running = len(os.popen('ps -ef | grep runn').readlines() > 2

if size_1 != size_2 or new_hash != current_hash or not running:
    for line in os.popen('killall discogssyncer').readlines():
        pass
    subprocess.Popen(['./discogssyncer', '--sync=false'])
