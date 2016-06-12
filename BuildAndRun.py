import os
import subprocess

# Update to the latest version
for line in os.popen('git fetch -p -q; git merge -q origin/master').readlines():
    print line.strip()
    
# Move the old version over
for line in os.popen('cp sync oldsync').readlines():
    print line.strip()

# Rebuild
for line in os.popen('go build ./...').readlines():
    print line.strip()

# Rebuild
for line in os.popen('go build').readlines():
    print line.strip()


size_1 = os.path.getsize('./oldsync')
size_2 = os.path.getsize('./sync')

if size_1 != size_2:
    for line in os.popen('killall sync').readlines():
        pass
    subprocess.Popen('./cardserver')
