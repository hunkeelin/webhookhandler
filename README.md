### INTRODUCTION ###
stable-1.0.3 added slack support. If you do not wnat slack support use stable-1.0.2

### KEY FEATURES ###
* Compare tags instead of doing an actual git pull
* Support concurrency, number of threads is controlled by a parameter in /etc/r9000k/r9000k.conf
* Supports lock and will wait 10 seconds if there's background process
* Cleans unwanted files and modules not define in Puppetfile

### Installation ###
- Built a .deb for debian/ubuntu and .rpm for redhats/centos
### BENCHMARK RESULTS ###
r10k takes 3 minutes and 10 seconds each run, r9000k takes no more than 10 seconds except the initial git clone. 

### Instructions ###
After installation of the package. Please edit your /etc/r9000k/r9000k.conf

- workdir: the working puppet directory. Most likely /etc/puppetlabs/code/environments/  PS: If you are to use a different directory for testing purposes, make sure that directory is created beforehand. 
- remote: the git repository that contains the Puppetfile. E.g remote="git@github.somecompany.com:puppet/environments.git"
- lockfile: the location of the lock file
- concur: the number of concurrency for performance purpose. Set it equal to the number of the threads for max performance. 
- slackurl: your slack bot url
- slackchannel: your slack channel e.g. "#general"

### Things that will break this program ###
- The puppetfile must be in the following syntax
```
// specifying tags
mod 'foo',
  :git => 'git@github.abc.com:puppet/foo.git',
  :tag => '1.1.0'

//specifying branch
mod 'user_management',
  :git => 'git@github.abc.com:puppet/user_management.git',
  :branch => 'dev'

//Default to master latest branch
mod 'account',
  :git => 'git@github.abc.com:puppet/puppet-accounts.git'
```
- Do not remove tag and recreate tags. r9000k will not update the repos if the tag are the same even if the contents are not. If you want to make a change make a new tag. 
- Your environment Puppetfile need to have a non master branch. R9000k will not setup an environment for master branch. 
- DO NOT touch /etc/r9000k/scripts I made those bash scripts because if I do it in go it's clobbering up my code and it makes it hard to read. 

### Change log ###
BUGFIX: Repository doesn't get updated even though describe tag shows the correct tag. 
BUGFIX: When there's no tag it will throw an error. Fixed by letting the loop continue. 
