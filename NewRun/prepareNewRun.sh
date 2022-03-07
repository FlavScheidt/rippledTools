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
LOGS_DIR="/var/log/rippled"
DB_DIR="/var/lib/rippled"

#I STRONGLY REMCOMEND MOVING THE LOGS AND DB TO A NEW LOCATION
#Dont put relative paths, they wont work
NEW_LOGS_DIR="/root/rippledLogs"
NEW_DB_DIR="/root/rippledDBs"

#Specifics for snts project
STDOUTLOGS_DIR="/root/sntrippled/my_build"
NEW_STDOUTLOGS_DIR="/root/sntrippled/my_build/logs"

GRPCLOGS_DIR="/root/sntrippled/grpc"
NEW_GRPCLOGS_DIR="/root/sntrippled/grpc/logs"

##########################
#	CHECK NEW DIRECTORIES
##########################
if [ ! -d "$NEW_LOGS_DIR" ]
then
	mkdir "$NEW_LOGS_DIR"
fi

if [ ! -d "$NEW_DB_DIR" ]
then
	mkdir "$NEW_DB_DIR"
fi

#Specifics for snt
if [ ! -d "$NEW_STDOUTLOGS_DIR" ]
then
	mkdir "$NEW_STDOUTLOGS_DIR"
fi

if [ ! -d "$NEW_GRPCLOGS_DIR" ]
then
	mkdir "$NEW_GRPCLOGS_DIR"
fi


##########################
#	START
##########################
readarray -t nodes < ./nodes.txt

for node in "${nodes}";
do

	#########################
	#	Rename logs
	#########################
	#Rename the last log with the timestamp a the begining of the log

	LOG_DATE=$(ssh ${node} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f1")
	LOG_HOUR=$(ssh ${node} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f1")
	LOG_MIN=$(ssh ${node} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f2")
	LOG_SEC=$(ssh ${node} "head -1 ${LOGS_DIR}/debug.log | cut -d ' ' -f2 | cut -d ':' -f3 | cut -d '.' -f1")


	ssh ${node} "mv ${LOGS_DIR}/debug.log ${NEW_LOGS_DIR}/debug_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.log"

	#########################
	#	Rename database
	#########################
	ssh ${node} "mkdir ${NEW_DB_DIR}/db_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}"
	ssh ${node} "mv ${DB_DIR}/db/* ${NEW_DB_DIR}/db_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}/"


	#Here down is for the specific use on Snts project
	########################
	# Rename stdout logs
	########################
	ssh ${node} "mv ${STDOUTLOGS_DIR}/log.out ${NEW_STDOUTLOGS_DIR}/log_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.out"

	########################
	# Rename GRPC logs
	########################
	ssh ${node} "mv ${GRPCLOGS_DIR}/log.out ${NEW_GRPCLOGS_DIR}/log_${LOG_DATE}_${LOG_HOUR}_${LOG_MIN}_${LOG_SEC}.out"

done