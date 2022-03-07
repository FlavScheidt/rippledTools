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
readarray -t nodes < ./nodes.txt

if [ "$1" == "report" ]
then
	###########################################################
	#	Report nodes status
	###########################################################


	for node in "${nodes}";
	do
		status=$(ssh ${node} "sntrippled/my_build/rippled server_info" | grep -q "active")

		if [ ${status} ]
		then
			echo "${node}		running"
		else
			status=$(ssh ${node} "sntrippled/my_build/rippled server_info" | grep -q "no responde from server")
			echo "${node}		stopped"
		fi
	done
fi