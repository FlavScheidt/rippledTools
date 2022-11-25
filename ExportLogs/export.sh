#!/bin/bash

##########################################################
#	This script is inteded to be used in one of the nodes
#	of the rippled cluster. I assume that /etc/hosts is
#	already configured with the nodes names and ips. I also
#	expect that ssh is configured to be used without password,
#	using only private/public key pairs authentication AND
#	that nodes have the same usernames.
#
#	Flaviene Scheidt de Cristo
# 	University of Luxembourg
#	May/2022
###########################################################

type=$1
echo "$type"

##########################
#	EDITABLE VARIABLES
##########################
DATA_DIR="/home/xrpl/data"
TARGET_DIR="~/data"
TEMP_DIR="/root/data"

TARGET_ADDRESS="flaviene.scheidt@sherlock.uni.lux"

LOGS_DIR="/root/var/log/rippled"
LOGS_GO_DIR="/root/gossipGoSnt"
LOGS_STDOUT_V_DIR="/root/rippled/my_build"
LOGS_STDOUT_G_DIR="/root/sntrippled/my_build"


##########################
#	LOAD NODE NAMES
##########################
readarray -t nodes < ./nodes.txt

##########################
#	Start
##########################

# if [ "$type" == "gossip" ]
# then
for n in "${nodes[@]}";
do
	echo "-----------${n}-----------"
	# ##########################
	# #	DIRECTORIES
	# ##########################
	if ssh ${n} "[ ! -d $DATA_DIR ]"
	then
		ssh ${n} "mkdir $DATA_DIR"
	fi

	TYPE_DIR="${DATA_DIR}/${type}"
	if ssh ${n} "[ ! -d $TYPE_DIR ]"
	then
		ssh ${n} "mkdir $TYPE_DIR"
	fi
	echo "Created directories on the source"

	rm -rf ${TYPE_DIR}/*
	ssh ${n} "rm -rf  $TYPE_DIR/*"
	echo "Directories are clean"

	if ssh ${TARGET_ADDRESS} "[ ! -d $TARGET_DIR ]"
	then
		ssh ${TARGET_ADDRESS} "mkdir $TARGET_DIR"
	fi

	if ssh ${TARGET_ADDRESS} "[ ! -d "$TARGET_DIR/${n}" ]"
	then
		ssh ${TARGET_ADDRESS} "mkdir $TARGET_DIR/${n}"
	fi

	TYPE_TARGET_DIR="${TARGET_DIR}/${n}/${type}"
	if ssh ${TARGET_ADDRESS} "[ ! -d $TYPE_TARGET_DIR ]"
	then
		ssh ${TARGET_ADDRESS} "mkdir $TYPE_TARGET_DIR"
	fi
	echo "Created directories on the target"

	ssh ${TARGET_ADDRESS} "rm -rf $TYPE_TARGET_DIR/*"
	echo "Directories are clean"

	ssh ${n} "cp ${LOGS_DIR}/debug.log ${TYPE_DIR}/debug.log"
	echo "Rippled logs"

	if [ "$type" == "vanillaGeneral" ] || [ "$type" == "vanillaValidator" ] || [ "$type" == "vanillaUNL" ]
	then
		ssh ${n} "cp ${LOGS_STDOUT_V_DIR}/log.out ${TYPE_DIR}/log_stdout.out"
		echo "Rippled stdout"
	else
		ssh ${n} "cp ${LOGS_STDOUT_G_DIR}/log.out ${TYPE_DIR}/log_stdout.out"
		echo "Stdout"

		ssh ${n} "cp ${LOGS_GO_DIR}/log.out ${TYPE_DIR}/log_go.out"
		echo "GossipSub stdout"

		# ssh ${n} "cp ${LOGS_GO_DIR}/trace.json ${TYPE_DIR}/trace.json"
		# echo "GossipSub trace"
	fi;

	if [ ! -d $TEMP_DIR ]
	then
		mkdir $TEMP_DIR
	fi
	echo "Created ${TEMP_DIR}"

	if [ ! -d "${TEMP_DIR}/${n}" ]
	then
		mkdir "${TEMP_DIR}/${n}"
	fi

	if [ ! -d "${TEMP_DIR}/${n}/${type}" ]
	then
		mkdir "${TEMP_DIR}/${n}/${type}"
	fi
	rm -rf ${TEMP_DIR}/${n}/${type}/*
	echo "Created ${TEMP_DIR}/${n}/${type}"


	scp -r ${n}:${TYPE_DIR}/* ${TEMP_DIR}/${n}/${type}
	scp -r ${TEMP_DIR}/${n}/${type}/* ${TARGET_ADDRESS}:${TYPE_TARGET_DIR}/
	
	rm -rf ${TEMP_DIR}/${n}/${type}
	echo "Logs moved to target"
done
	
	# fi}