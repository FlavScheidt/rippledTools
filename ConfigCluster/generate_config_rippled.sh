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
rippledmon_IP="192.168.20.52"
#CONFIG_DIR is the folder where you want to put the files on the target servers
CONFIG_DIR="/root/config/"

##########################
#	Craeting some directories
#		(if needed)
##########################
#if [ ! -d "$CONFIG_DIR" ]
#then
#	mkdir "$CONFIG_DIR"
#fi

##########################
#	Are we cleaning?
# 	If not, generate files
##########################
if [ "$1" == "clean" ]
then
	rm -rf validators_* rippled_*
else

	##########################
	#	GET LIST OF NODES
	##########################
	nodes=($(cat ClusterConfigSmall.csv | cut -d ',' -f2))
	# nodes="lotus"
	# echo ${nodes[@]}

	#########################
	#	Generate config files
	#########################
	#Iterate over list of nodes
	for n in "${nodes[@]}";
	do

		#Nothing to see here, just being sure the config directory exists
		if ssh ${n} "[ ! -d ${CONFIG_DIR} ]"
		then
			ssh ${n} "mkdir ${CONFIG_DIR}"
		fi

		# RYCB Uncomment if you want to use the rippled.cfg that is already configured on your node
		# scp ${n}:/opt/local/etc/rippled.cfg ./rippled_${n}.cfg

		# RYCB Uncomment if you want to use the local rippled.cfg file (you need to provide the keys, then)
		cp rippled.cfg rippled_${n}.cfg

		# Create validators file
		cp validators.txt validators_${n}.txt
		echo "" | tee -a validators_${n}.txt >/dev/null

		# Add key
		# RYCB comment if you already configured the key, if not, be sure to put it on the keys folder
		echo "" | tee -a rippled_${n}.cfg >/dev/null
		cat ./keys/${n}.txt | tee -a rippled_${n}.cfg >/dev/null
		echo "" | tee -a rippled_${n}.cfg >/dev/null

		# Insight
		# RYCB If you don't want to configure rippledmon, comment this line
		# echo "[insight]" | tee -a rippled_${n}.cfg >/dev/null
		# echo "server=statsd" | tee -a rippled_${n}.cfg >/dev/null
		# echo "addres=${rippledmon_IP}:8125" | tee -a rippled_${n}.cfg >/dev/null
		# echo "prefix=${n}_" | tee -a rippled_${n}.cfg >/dev/null
		# echo "" | tee -a rippled_${n}.cfg >/dev/null


		# Print the IPS of the UNL
		# First we need to get the ips
		readarray -t unl < ./unl/fullySmall/${n}.txt
		# echo ${unl[@]}

		#Print header into rippled.cfg
		echo "" | tee -a rippled_${n}.cfg >/dev/null
		echo "[ips_fixed]" | tee -a rippled_${n}.cfg >/dev/null

		#iterate over unl and print ips into rippled.cfg and keys into validators.txt
		for peer in "${unl[@]}";
		do
			# Get IP and print to rippled.cfg
			ip=($(cat ClusterConfig.csv | grep -i ",${peer}," | cut -d "," -f1))
			echo "$ip 51235" | tee -a rippled_${n}.cfg >/dev/null

			#Get key and print to validators.txt
			key=($(cat ClusterConfig.csv | grep -i ",${peer}," | cut -d "," -f3))
			echo "$key" | tee -a validators_${n}.txt >/dev/null

		done
		echo "" | tee -a rippled_${n}.cfg >/dev/null

		# Send to server
		scp ./rippled_${n}.cfg ${n}:${CONFIG_DIR}/rippled.cfg
		scp ./validators_${n}.txt ${n}:${CONFIG_DIR}/validators.txt

	 done
fi
