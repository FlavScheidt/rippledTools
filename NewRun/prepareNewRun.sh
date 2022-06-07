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

##########################
#	EDITABLE VARIABLES
##########################
LOGS_DIR="/root/var/log/rippled"
DB_DIR="/var/lib/rippled"

#I STRONGLY REMCOMEND MOVING THE LOGS AND DB TO A NEW LOCATION
#Dont put relative paths, they wont work
NEW_LOGS_DIR="/home/xrpl/rippledLogs"
NEW_DB_DIR="/home/xrpl/rippledDBs"

#Specifics for snts project
STDOUTLOGS_DIR="/root/sntrippled/my_build"
NEW_STDOUTLOGS_DIR="/home/xrpl/rippledStdOutLogs"

STDOUTVANILLALOGS_DIR="/root/rippled/my_build"
NEW_STDOUTVANILLALOGS_DIR="/home/xrpl/rippledVanillaStdOutLogs"

GRPCLOGS_DIR="/root/gossipGoSnt"
NEW_GRPCLOGS_DIR="/home/xrpl/rippledgRPCLogs"


##########################
#	START
##########################
readarray -t nodes < ./nodes.txt

for n in "${nodes[@]}";
do
	echo ${n}

	# ##########################
	# #	CHECK NEW DIRECTORIES
	# ##########################
	# if ssh ${n} "[ ! -d $NEW_LOGS_DIR ]"
	# then
	# 	ssh ${n} "mkdir $NEW_LOGS_DIR"
	# fi

	# if ssh ${n} "[ ! -d $NEW_DB_DIR ]"
	# then
	# 	ssh ${n} "mkdir $NEW_DB_DIR"
	# fi

	# #Specifics for snt
	# if ssh ${n} "[ ! -d $NEW_STDOUTLOGS_DIR ]"
	# then
	# 	ssh ${n} "mkdir $NEW_STDOUTLOGS_DIR"
	# fi

	# if ssh ${n} "[ ! -d $NEW_STDOUTVANILLALOGS_DIR ]"
	# then
	# 	ssh ${n} "mkdir $NEW_STDOUTVANILLALOGS_DIR"
	# fi

	# if ssh ${n} "[ ! -d $NEW_GRPCLOGS_DIR ]"
	# then
	# 	ssh ${n} "mkdir $NEW_GRPCLOGS_DIR"
	# fi


	# if ssh ${n} "[ -f ${LOGS_DIR}/debug.log ]"
	# then

	# 	#########################
	# 	#	Rename logs
	# 	#########################
	# 	#Rename the last log with the timestamp a the begining of the log

	# 	LOG_DATE=$(ssh ${n} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f1")
	# 	LOG_HOUR=$(ssh ${n} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f1")
	# 	LOG_MIN=$(ssh ${n} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f2")
	# 	LOG_SEC=$(ssh ${n} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f3 | cut -d '.' -f1")


	# 	ssh ${n} "mv ${LOGS_DIR}/debug.log ${NEW_LOGS_DIR}/debug_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.log"

	# 	#########################
	# 	#	Rename database
	# 	#########################
	# 	if ssh ${n} "[ -d ${DB_DIR}/db ]"
	# 	then
	# 		ssh ${n} "mkdir ${NEW_DB_DIR}/db_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}"
	# 		ssh ${n} "mv ${DB_DIR}/db/* ${NEW_DB_DIR}/db_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}/"
	# 	else 
	# 		echo "No database"
	# 	fi

	# 	#Here down is for the specific use on Snts project
	# 	########################
	# 	# Rename stdout logs
	# 	########################
	# 	if ssh ${n} "[ -f ${STDOUTLOGS_DIR}/log.out ]"
	# 	then
	# 		ssh ${n} "mv ${STDOUTLOGS_DIR}/log.out ${NEW_STDOUTLOGS_DIR}/log_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.out"
	# 	else 
	# 		echo "No stdout log"
	# 	fi

	# 	########################
	# 	# Rename stdout logs
	# 	########################
	# 	if ssh ${n} "[ -f ${STDOUTVANILLALOGS_DIR}/log.out ]"
	# 	then
	# 		ssh ${n} "mv ${STDOUTVANILLALOGS_DIR}/log.out ${NEW_STDOUTVANILLALOGS_DIR}/log_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.out"
	# 	else 
	# 		echo "No stdout vanilla log"
	# 	fi

	# 	########################
	# 	# Rename GRPC logs
	# 	########################
	# 	if ssh ${n} "[ -f ${GRPCLOGS_DIR}/log.out ]"
	# 	then
	# 		ssh ${n} "mv ${GRPCLOGS_DIR}/log.out ${NEW_GRPCLOGS_DIR}/log_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.out"
	# 	else 
	# 		echo "No gRPC logs"
	# 	fi
	# else 
		echo "No log to rename"
		echo "Deleting databases and cleaning directories..."

		ssh ${n} "rm -rf ${DB_DIR}/db/*"
		ssh ${n} "rm -rf /var/log/ripple/*"
		ssh ${n} "rm -rf ${LOGS_DIR}/*"
		ssh ${n} "rm -rf ${STDOUTLOGS_DIR}/log.out"
		ssh ${n} "rm -rf ${STDOUTVANILLALOGS_DIR}/log.out"
		ssh ${n} "rm -rf ${GRPCLOGS_DIR}/log.out"
		ssh ${n} "rm -rf ${GRPCLOGS_DIR}/trace.json"
	# fi

done