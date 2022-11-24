from pssh.clients import ParallelSSHClient

from pssh.utils import enable_host_logger

import csv
from subprocess import call
import time



##########################
#   EDITABLE VARIABLES
##########################
PATH="/home/xrpl/rippledTools/"
RIPPLED_PATH="./sntrippled/my_build/"
RIPPLED_CONFIG="/root/config/rippled.cfg"
RIPPLED_QUORUM="15"

GOSSIPSUB_PATH="/root/gossipGoSnt/"

GOPATH="/usr/local/go/bin/"

#########################
#   Clean envinroments for new execution
#########################
rc = call(PATH+"/NewRun/prepareNewRun.sh", shell=True)
# print(rc)

#########################
#    Nodes names
#########################
hosts=[]
with open(PATH+'/ConfigCluster/ClusterConfigSmall.csv', newline='') as file:
    reader = csv.reader(file, delimiter=',')
    for row in reader:
        hosts.append(row[1])

clientRippled = ParallelSSHClient(hosts)
clientGossip = ParallelSSHClient(hosts)

clientRippledStop = ParallelSSHClient(hosts)
clientGossipStop = ParallelSSHClient(hosts)

# enable_host_logger()

#########################
#    Load Parameters
#       Form: -d, -dlo, -dhi, -dscore, -dlazy, -dout, -gossipFactor, -InitialDelay, -Interval
#########################
parameters=[]
with open(PATH+'/Orchestrator/parameters.csv', newline='') as file:
    reader = csv.reader(file, delimiter=',')
    for row in reader:
        parameters.append(row)

#########################
#    Iterate over list of parameters
#########################
rippled = RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" --quorum "+RIPPLED_QUORUM
rippledStop = RIPPLED_PATH+"rippled --conf="+RIPPLED_CONFIG+" stop"
gossip_kill = "pkill -f 'gossipGoSnt'"

experiment="unl"
for param in parameters:

    gossipsub = "cd "+GOSSIPSUB_PATH+" && "+GOPATH+"/go run . -type="+experiment+" -d="+param[0]+" -dlo="+param[1]+" -dhi="+param[2]+" -dscore="+param[3]+" -dlazy="+param[4]+" -dout="+param[5]

    ########### Run
    output_rippled = clientRippled.run_command(rippled, use_pty=True)
    time.sleep(15)
    output_gossipsub = clientGossip.run_command(gossipsub, use_pty=True)

    # client.join(output_rippled, timeout=60)#900
    # client.join(output_gossipsub)

    # time.sleep(30)#900
    # print("stopping...")

    ########## Stop
    # output_rippledStop = clientRippledStop.run_command(rippledStop, use_pty=True)

    #Need to get the pid of the gossipsub process to send a kill sign to the node
    # output_gossipkill = clientGossipStop.run_command(gossip_kill, use_pty=True)


    clientRippled.join(output_rippled, timeout=60)#consume_output=True)
    clientGossip.join(output_gossipsub, timeout=60)

    # for host_output in output_gossipkill:
    #     hostname = host_output.host
    #     # stdout = list(host_output.stdout)
    #     print("Gossipkill Host %s: exit code %s" % (
    #           hostname, host_output.exit_code))#, stdout))

    for host_output in output_rippled:
        hostname = host_output.host
        stdout = list(host_output.stdout)
        print("Rippled Host %s: exit code %s" % (
              hostname, host_output.exit_code))
        # if host_output.exit_code != 0:
        stdout = list(host_output.stdout)
        print(stdout)

    # for host_output in output_rippledStop:
    #     hostname = host_output.host
    #     stdout = list(host_output.stdout)
    #     print("Rippled stop Host %s: exit code %s, output %s" % (
    #           hostname, host_output.exit_code, stdout))

    # for host_output in output_gossipsub:
    #     hostname = host_output.host
    #     # stdout = list(host_output.stdout)
    #     print("Gossipsub Host %s: exit code %s" % (
    #           hostname, host_output.exit_code))#, stdout))
    #     if host_output.exit_code != 0:
    #         stdout = list(host_output.stdout)
    #         print(stdout)
