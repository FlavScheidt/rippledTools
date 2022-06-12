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

arr=("time 	node 	closed	index 	time 	parent 	openIndex")
##########################
#	START
##########################
readarray -t nodes < ./nodes.txt

if [ "$1" == "vanilla" ]
then
	dir="rippled"
else
	dir="sntrippled"
fi

for n in "${nodes[@]}";
do
	# echo "------------------------------------------"
	# echo ${n}
	ssh ${n} "./${dir}/my_build/rippled ledger" >> tmp_${n}.out
	# cat tmp_${n}.out
done

for n in "${nodes[@]}";
do
	closed=$(cat tmp_${n}.out | grep -m 1 "\"hash\"" | cut -d ':' -f2 | cut -d '"' -f2)
	parent=$(cat tmp_${n}.out | grep -m 1 "\"parent_hash\"" | cut -d ':' -f2 | cut -d '"' -f2)
	index=$(cat tmp_${n}.out | grep -m 1 "\"ledger_index\"" | cut -d ':' -f2 | cut -d '"' -f2)
	time=$(cat tmp_${n}.out | grep -m 1 "\"close_time\"" | cut -d ':' -f2 | cut -d '"' -f2)
	openIndex=$(cat tmp_${n}.out | grep "\"ledger_index\"" | sed -n 2p | cut -d ':' -f2 | cut -d '"' -f2)

	arr+=("$(date +%T) ${n}	${closed} 	${index} 	${time} 	${parent} 	${openIndex}")
	rm tmp_${n}.out
done

for line in "${arr[@]}";
do
	echo $line
done