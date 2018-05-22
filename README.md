### INTRODUCTION ###
Jenkins written in go. Smaller memory footprint. Take less resource and many time faster. It is also more secure by default. 


### Please edit /etc/genkins/genkins.conf ###
- concur: The number of thread use. Set to the number of CPU for max performance
- lockfile: the location of the lockfile
- certpath: the webcert path 
- keypath: the webkey path 
- bindaddr: bind to the interface you want the service to run. If left blank it will bind to all interface
- port: binding port
- apikey: ignore this, it's for future use. 
- hosts: ignore this for now, it's for future use. 
- jobdir: the job directory for all the genkins jobs. Will explain in the following session. 

### Creating jobs for hooks ###
Assume the following:
- jobdir = /tmp/jobs
- webhook comes from https://github.somecompany.com/IT/testrepo

- create the directory /tmp/jobs/github.somecompany.com/IT/testrepo
- create the config file /tmp/jobs/github.somecompany.com/IT/testrepo/config with the following content

#### /tmp/jobs/github.somecompany.com/IT/testrepo/config ####
secret = mysecret // this secret is the same secret you setup in your githooks via github. This is to verify the signiture to prevent dubious server from sending request to genkins. 

- for any jobs create any files ending with *.conf in the directory e.g /tmp/jobs/github.somecompany.com/IT/testrepo/foo.conf

#### /tmp/jobs/github.somecompany.com/IT/testrepo/foo.conf ####

run=/tmp/foo.sh


- The above example is bascially equivalent to jenkins's execute shell feature. Therefore genkins should never be ran as root. 
