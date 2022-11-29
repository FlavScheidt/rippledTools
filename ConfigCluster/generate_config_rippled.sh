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
rippledmon_IP="192.168.20.58"
#CONFIG_DIR is the folder where you want to put the files on the target servers
CONFIG_DIR="/root/config/"

UNL_DIR="/root/gossipGoSnt/clusterConfig"

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
	if [ "$1" == "small" ]
	then
		nodes=($(cat ClusterConfigSmall.csv | cut -d ',' -f2))
	else
		nodes=($(cat ClusterConfig.csv | cut -d ',' -f2))
	fi
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

		# Insight
		# RYCB If you don't want to configure rippledmon, comment this line
		# if [ ${n} = 'lotus' ] 
		# then
		# 	echo "" | tee -a rippled_${n}.cfg >/dev/null
		# 	echo "[insight]" | tee -a rippled_${n}.cfg >/dev/null
		# 	echo "server=statsd" | tee -a rippled_${n}.cfg >/dev/null
		# 	echo "address=${rippledmon_IP}:8125" | tee -a rippled_${n}.cfg >/dev/null
		# 	echo "prefix=${n}" | tee -a rippled_${n}.cfg >/dev/null
		# 	# echo "" | tee -a rippled_${n}.cfg >/dev/null
		# fi

		# Add key
		# RYCB comment if you already configured the key, if not, be sure to put it on the keys folder
		echo "" | tee -a rippled_${n}.cfg >/dev/null
		cat ./keys/${n}.txt | tee -a rippled_${n}.cfg >/dev/null
		echo "" | tee -a rippled_${n}.cfg >/dev/null

		nodeIP=($(cat ClusterConfig.csv | grep -i ",${n}," | cut -d "," -f1))
		nodeKey=($(cat ClusterConfig.csv | grep -i ",${n}," | cut -d "," -f3))

		# Print the IPS of the UNL
		# First we need to get the ips
		if [ "$1" == "unl" ]
		then
			unl=($(cat ClusterConfig.csv | grep -i ",${n}," | cut -d "," -f4))
			readarray -t unl < ${UNL_DIR}/unl/${unl}.txt
		elif [ "$1" == "general" ]
		then
			readarray -t unl < ./unl/fullyConnected/${n}.txt
		elif [ "$1" == "validator" ]
		then
			readarray -t unl < ${UNL_DIR}/validator/${n}.txt
		elif [ "$1" == "small" ]
		then
			readarray -t unl < /root/rippledTools/ConfigCluster/unl/fullySmall/${n}.txt
		fi
		# echo ${unl[@]}

		#Print header into rippled.cfg
		echo "" | tee -a rippled_${n}.cfg >/dev/null
		echo "[ips_fixed]" | tee -a rippled_${n}.cfg >/dev/null

		# #iterate over unl and print ips into rippled.cfg and keys into validators.txt
		for peer in "${unl[@]}";
		do
			if [ "$[peer]" != "${n}" ] 
			then
				# Get IP and print to rippled.cfg
				ip=($(cat ClusterConfig.csv | grep -i ",${peer}," | cut -d "," -f1))
				echo "$ip 51235" | tee -a rippled_${n}.cfg >/dev/null

				#Get key and print to validators.txt
				key=($(cat ClusterConfig.csv | grep -i ",${peer}," | cut -d "," -f3))
				echo "$key" | tee -a validators_${n}.txt >/dev/null
			fi
		done

		#Now we need to insert the ips and keys of the guys we are SENDINg messages to (only for the perUNL)
		if [ "$1" == "unl" ]
		then
			#Iterate over the UNL files to see in which unl the node is present
			FILES="${UNL_DIR}/unl/*"
			for f in $FILES
			do
				if grep -Fxq "${n}" ${f}
				then
					unlName=$(echo ${f} | cut -d "." -f1 | cut -d "/" -f6)

					#Go to the config file and get the keys from the nodes that have this unl
					ips=($(cat ClusterConfig.csv | grep -i ",${unlName}" | cut -d "," -f1))
					keys=($(cat ClusterConfig.csv | grep -i ",${unlName}" | cut -d "," -f3))


					for ip in "${ips[@]}"
					do 
						isInFile=$(cat rippled_${n}.cfg | grep -c ${ip})
						if [ $isInFile -eq 0 ] && [ "${ip}" != "${nodeIP}" ]
						then
							echo "${ip} 51235" |  tee -a rippled_${n}.cfg >/dev/null;
						fi
					done
					for key in "${keys[@]}"
					do 						
						isInFile=$(cat validators_${n}.txt | grep -c ${key})
						if [ $isInFile -eq 0 ] && [ "${key}" != "${nodeKey}" ]
						then
							echo "${key}" | tee -a validators_${n}.txt >/dev/null;
						fi
					done  
				fi
			done
		fi

		echo "" | tee -a rippled_${n}.cfg >/dev/null

		# Send to server
		scp ./rippled_${n}.cfg ${n}:${CONFIG_DIR}/rippled.cfg
		scp ./validators_${n}.txt ${n}:${CONFIG_DIR}/validators.txt

	 done
fi
